package tree

import (
	"fmt"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
)

type (
	definition struct {
		defaultValue string
		records      map[string]*record //maps label_expression to FeatureOverrideLookup
	}
	record struct {
		overrideName string
		value        string
		priority     int
		exists       bool
		records      map[string]*record //maps next in chain label_expression to inner overrides FeatureOverrideLookup
	}
	InMemoryTreeStorage struct {
		sync.RWMutex
		tree map[string]*definition
	}
)

func NewInMemoryStorage() storage.Storage {
	return &InMemoryTreeStorage{tree: map[string]*definition{}}
}

func (s *InMemoryTreeStorage) Define(ffd storage.FeatureFlagDefinition) {
	s.Lock()
	defer s.Unlock()
	if def, ok := s.tree[ffd.FeatureName]; ok {
		logrus.WithField("FeatureName", ffd.FeatureName).WithField("NewDefaultValue", ffd.DefaultValue).Error("Definition already exists, defaultValue will be overridden")
		def.defaultValue = ffd.DefaultValue
		return
	}
	s.tree[ffd.FeatureName] = &definition{
		defaultValue: ffd.DefaultValue,
		records:      make(map[string]*record),
	}
}

func (s *InMemoryTreeStorage) Override(ffo storage.FeatureFlagOverride) {
	s.Lock()
	defer s.Unlock()
	if d, ok := s.tree[ffo.FeatureName]; ok {
		descriptors := labelsToDescriptors(ffo.Labels)
		if duplicate := d.findDuplicate(descriptors); duplicate != nil {
			logrus.WithField("FeatureName", ffo.FeatureName).Errorf("Duplicate found %#v", duplicate)
		}
		d.insertRecord(ffo.Value, ffo.Origin, ffo.Priority, descriptors)
		return
	}
	logrus.WithField("FeatureName", ffo.FeatureName).WithField("Action", "Override").Errorf("Definition not found")
}

func (s *InMemoryTreeStorage) Find(featureName string, labels map[string]string) *pb.FeatureFlag {
	s.RLock()
	defer s.RUnlock()
	if d, ok := s.tree[featureName]; ok {
		descriptors := labelsToDescriptors(labels)
		if foundRecords, _ := matchLabels(d.records, descriptors); len(foundRecords) > 0 {
			sort.Slice(foundRecords, func(i, j int) bool {
				return foundRecords[i].priority > foundRecords[j].priority
			})
			foundRecord := foundRecords[0]
			if foundRecord.exists {
				return &pb.FeatureFlag{Origin: foundRecord.overrideName, FeatureName: featureName, Value: foundRecord.value}
			}
		}
		return &pb.FeatureFlag{FeatureName: featureName, Value: d.defaultValue}
	}
	return nil
}
func (d *definition) insertRecord(value, overrideName string, priority int, descriptors []string) {
	r := &record{value: value, overrideName: overrideName, priority: priority, records: map[string]*record{}, exists: true}
	insert(d.records, descriptors, r)
}
func insert(records map[string]*record, descriptors []string, r *record) {
	length := len(descriptors)
	switch length {
	case 0:
		logrus.Error("It's impossible. Something went wrong")
	case 1:
		records[descriptors[0]] = r
	default:
		nestedRecord, ok := records[descriptors[0]]
		if !ok {
			nestedRecord = &record{records: map[string]*record{}, exists: false}
			records[descriptors[0]] = nestedRecord
		}
		insert(nestedRecord.records, descriptors[1:], r)
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
	removeRecord(d.records, descriptors)
}

func removeRecord(records map[string]*record, descriptors []string) {
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
			}
		}
	default:
		nestedRecord, ok := records[descriptors[0]]
		if ok {
			removeRecord(nestedRecord.records, descriptors[1:])
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
	if foundRecord, ok := records[descriptors[0]]; ok {
		foundRecords = append(foundRecords, foundRecord)
		nestedFound, nestedCompleted := matchLabels(foundRecord.records, descriptors[1:])
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

func (s *InMemoryTreeStorage) RemoveDefinition(featureName string) {
	s.Lock()
	defer s.Unlock()
	delete(s.tree, featureName)
}

func (s *InMemoryTreeStorage) RemoveOverride(featureName string, labels map[string]string) {
	s.Lock()
	defer s.Unlock()
	if d, ok := s.tree[featureName]; ok {
		descriptors := labelsToDescriptors(labels)
		d.removeRecord(descriptors)
		return
	}
	logrus.WithField("FeatureName", featureName).WithField("Action", "RemoveOverride").Error("Definition not found")
}
