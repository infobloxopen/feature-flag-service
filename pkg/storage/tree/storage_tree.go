package tree

import (
	"fmt"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"

	// "github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"

	ffv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
)

type (
	definition struct {
		defaultValue string
		records      map[string]*record // maps label_expression to FeatureOverrideLookup
	}
	record struct {
		overrideName string
		value        string
		priority     int
		exists       bool
		depth        int
		records      map[string]*record // maps next in chain label_expression to inner overrides FeatureOverrideLookup
	}
	InMemoryTreeStorage struct {
		sync.RWMutex
		tree map[string]*definition
	}
)

// NewInMemoryStorage ...
func NewInMemoryStorage() storage.Storage {
	return &InMemoryTreeStorage{tree: map[string]*definition{}}
}

// Define ...
func (s *InMemoryTreeStorage) Define(obj *ffv1.FeatureFlag) {
	s.Lock()
	defer s.Unlock()
	logrus.WithField("FeatureName", obj.Spec.FeatureID).Trace("Defining Feature")
	if def, ok := s.tree[obj.Spec.FeatureID]; ok {
		logrus.WithField("FeatureName", obj.Spec.FeatureID).WithField("NewDefaultValue", obj.Spec.Value).Error("Definition already exists, defaultValue will be overridden")
		def.defaultValue = obj.Spec.Value
		return
	}
	s.tree[obj.Spec.FeatureID] = &definition{
		defaultValue: obj.Spec.Value,
		records:      make(map[string]*record),
	}
}

// Override ...
func (s *InMemoryTreeStorage) Override(obj *ffv1.FeatureFlagOverride) {
	s.Lock()
	defer s.Unlock()
	logrus.WithField("FeatureName", obj.Spec.FeatureID).Trace("Defining FeatureFlagOverride")
	if d, ok := s.tree[obj.Spec.FeatureID]; ok {
		descriptors := labelsToDescriptors(obj.Labels())
		if duplicate := d.findDuplicate(descriptors); duplicate != nil {
			logrus.WithField("FeatureName", obj.Spec.FeatureID).Errorf("Duplicate found %#v", duplicate)
		}
		d.insertRecord(obj.Spec.Value, obj.Spec.OverrideName, obj.Spec.Priority, descriptors)
		return
	}
	logrus.WithField("FeatureName", obj.Spec.FeatureID).WithField("Action", "Override").Errorf("Definition not found")
}

// FindAll ...
func (s *InMemoryTreeStorage) FindAll(labels map[string]string) []*pb.FeatureFlag {
	s.RLock()
	defer s.RUnlock()
	var resultSet []*pb.FeatureFlag
	for featureName, d := range s.tree {
		descriptors := labelsToDescriptors(labels)
		if r := d.findByDescriptors(descriptors); r != nil {
			resultSet = append(resultSet, &pb.FeatureFlag{Origin: r.overrideName, FeatureName: featureName, Value: r.value})
			continue
		}
		resultSet = append(resultSet, &pb.FeatureFlag{FeatureName: featureName, Value: d.defaultValue})
	}
	return resultSet
}

// Find ...
func (s *InMemoryTreeStorage) Find(featureName string, labels map[string]string) *pb.FeatureFlag {
	s.RLock()
	defer s.RUnlock()
	logrus.WithField("FeatureName", featureName).Trace("Finding Feature")
	if d, ok := s.tree[featureName]; ok {
		descriptors := labelsToDescriptors(labels)
		if r := d.findByDescriptors(descriptors); r != nil {
			return &pb.FeatureFlag{Origin: r.overrideName, FeatureName: featureName, Value: r.value}
		}
		return &pb.FeatureFlag{FeatureName: featureName, Value: d.defaultValue}
	}
	return nil
}

func (d *definition) findByDescriptors(descriptors []string) *record {
	if foundRecords, _ := matchLabels(d.records, descriptors); len(foundRecords) > 0 {
		names := []string{}
		for _, r := range foundRecords {
			names = append(names, fmt.Sprintf("%d %d- %s (%s)", r.depth, r.priority, r.overrideName, r.value))
		}
		logrus.WithField("Before", names).Trace("findByDescriptors")
		sort.Slice(foundRecords, func(i, j int) bool {
			if foundRecords[i].priority > foundRecords[j].priority {
				return true
			}
			if foundRecords[i].priority < foundRecords[j].priority {
				return false
			}
			return foundRecords[i].depth > foundRecords[j].depth
		})
		names = []string{}
		for _, r := range foundRecords {
			names = append(names, fmt.Sprintf("%d %d- %s (%s)", r.depth, r.priority, r.overrideName, r.value))
		}
		logrus.WithField("After", names).Trace("findByDescriptors")
		return foundRecords[0]
	}
	return nil
}

func (d *definition) insertRecord(value, overrideName string, priority int, descriptors []string) {
	r := &record{value: value, overrideName: overrideName, priority: priority, records: map[string]*record{}, exists: true}
	if len(d.records) > 0 {
		r.records = d.records
	}
	insert(d.records, descriptors, r, 0)
}

func insert(records map[string]*record, descriptors []string, r *record, depth int) {
	length := len(descriptors)
	switch length {
	case 0:
		logrus.Error("It's impossible. Something went wrong")
	case 1:
		r.depth = depth
		records[descriptors[0]] = r
		logrus.WithField("depth", depth).WithField("descriptor", descriptors[0]).WithField("value", r.value).Trace("insert")
	default:
		nestedRecord, ok := records[descriptors[0]]
		if !ok {
			nestedRecord = &record{records: map[string]*record{}, exists: false, depth: depth}
			logrus.WithField("depth", depth).WithField("descriptor", descriptors[0]).Trace("insert")
			records[descriptors[0]] = nestedRecord
		} else {
			logrus.WithField("depth", depth).WithField("descriptor", descriptors[0]).Trace("existing")
		}
		insert(nestedRecord.records, descriptors[1:], r, depth+1)
	}
}

func (d *definition) findDuplicate(descriptors []string) *record {
	foundRecords, completed := matchLabels(d.records, descriptors)
	if completed {
		return foundRecords[len(foundRecords)-1]
	}
	return nil
}

func (d *definition) removeRecord(descriptors []string) {
	removeRecord(d.defaultValue, d.records, descriptors)
}

func removeRecord(defaultValue string, records map[string]*record, descriptors []string) {
	length := len(descriptors)
	switch length {
	case 0:
		logrus.Error("It's impossible. Something went wrong")
	case 1:
		key := descriptors[0]
		if record, ok := records[key]; ok {
			if len(record.records) == 0 {
				delete(records, key)
			} else {
				record.clearValue()
				record.value = defaultValue
			}
		}
	default:
		nestedRecord, ok := records[descriptors[0]]
		if ok {
			removeRecord(defaultValue, nestedRecord.records, descriptors[1:])
		}
	}
}

func (r *record) clearValue() {
	r.exists = false
	r.priority = 0
	r.overrideName = ""
	r.value = ""
}

func matchLabels(records map[string]*record, descriptors []string) (foundRecords []*record, completed bool) {
	if len(descriptors) == 0 {
		completed = true
		return
	}
	var i int
	for _, descriptor := range descriptors {
		if _, ok := records[descriptor]; ok {
			break
		}
		i++
	}
	if i >= len(descriptors) {
		return
	}
	if foundRecord, ok := records[descriptors[i]]; ok {
		foundRecords = append(foundRecords, foundRecord)
		nestedFound, nestedCompleted := matchLabels(foundRecord.records, descriptors[i+1:])
		completed = nestedCompleted
		foundRecords = append(foundRecords, nestedFound...)
	}
	return
}

func labelsToDescriptors(labels map[string]string) []string {
	var descriptors []string
	for k, v := range labels {
		descriptors = append(descriptors, fmt.Sprintf("%q=%q", k, v))
	}
	sort.Strings(descriptors)
	return descriptors
}

// RemoveDefinition ...
func (s *InMemoryTreeStorage) RemoveDefinition(obj *ffv1.FeatureFlag) {
	s.Lock()
	defer s.Unlock()
	delete(s.tree, obj.Spec.FeatureID)
}

// RemoveOverride ...
func (s *InMemoryTreeStorage) RemoveOverride(obj *ffv1.FeatureFlagOverride) {
	s.Lock()
	defer s.Unlock()
	if d, ok := s.tree[obj.Spec.FeatureID]; ok {
		descriptors := labelsToDescriptors(obj.Spec.LabelSelector.MatchLabels)
		d.removeRecord(descriptors)
		return
	}
	logrus.WithField("FeatureName", obj.Spec.FeatureID).WithField("Action", "RemoveOverride").Error("Definition not found")
}
