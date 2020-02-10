package flat

import (
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
)

type (
	record struct {
		priority     int
		overrideName string
		labels       map[string]string
		value        string
	}
	definition struct {
		defaultValue string
		records      []*record
	}
	InMemoryFlatStorage struct {
		sync.RWMutex
		table map[string]*definition
	}
)

func NewInMemoryStorage() storage.Storage {
	return &InMemoryFlatStorage{table: map[string]*definition{}}
}

func (s *InMemoryFlatStorage) Define(ffd storage.FeatureFlagDefinition) {
	s.Lock()
	defer s.Unlock()
	if def, ok := s.table[ffd.FeatureName]; ok {
		logrus.WithField("FeatureName", ffd.FeatureName).WithField("NewDefaultValue", ffd.DefaultValue).Error("Definition already exists, defaultValue will be overridden")
		def.defaultValue = ffd.DefaultValue
		return
	}
	s.table[ffd.FeatureName] = &definition{
		defaultValue: ffd.DefaultValue,
		records:      []*record{},
	}
}

func (s *InMemoryFlatStorage) Override(ffo storage.FeatureFlagOverride) {
	s.Lock()
	defer s.Unlock()
	if d, ok := s.table[ffo.FeatureName]; ok {
		if duplicates := d.findDuplicate(ffo.Labels); len(duplicates) > 0 {
			logrus.WithField("FeatureName", ffo.FeatureName).Errorf("Duplicates found %#v", duplicates)
		}
		d.insertRecord(ffo.Value, ffo.Origin, ffo.Labels, ffo.Priority)
		return
	}
	logrus.WithField("FeatureName", ffo.FeatureName).WithField("Action", "Override").Error("Definition not found")
}

func (s *InMemoryFlatStorage) Find(featureName string, labels map[string]string) *pb.FeatureFlag {
	s.RLock()
	defer s.RUnlock()
	if d, ok := s.table[featureName]; ok {
		result := &pb.FeatureFlag{FeatureName: featureName, Value: d.defaultValue}
		for _, r := range d.records {
			if r.matchLabels(labels) {
				result.Value = r.value
				result.Origin = r.overrideName
				break
			}
		}
		return result
	}
	return nil
}

func (d *definition) findDuplicate(labels map[string]string) []*record {
	var duplicates []*record
	for _, r := range d.records {
		if r.matchLabelsOne2One(labels) {
			duplicates = append(duplicates, r)
		}
	}
	return duplicates
}

func (d *definition) insertRecord(value string, overrideName string, labels map[string]string, priority int) {
	record := &record{
		priority:     priority,
		overrideName: overrideName,
		labels:       labels,
		value:        value,
	}
	records := d.records
	index := sort.Search(len(records), func(i int) bool { return records[i].priority < record.priority })
	records = append(records, nil)
	copy(records[index+1:], records[index:])
	records[index] = record
	d.records = records
}

func (d *definition) removeRecords(labels map[string]string) {
	var otherRecords []*record
	for _, r := range d.records {
		if !r.matchLabelsOne2One(labels) {
			otherRecords = append(otherRecords, r)
		}
	}
	d.records = otherRecords
}

func (r *record) matchLabels(labels map[string]string) bool {
	if len(r.labels) > len(labels) {
		return false
	}
	return findAllLabelsInLabels(r.labels, labels)
}

func (r *record) matchLabelsOne2One(labels map[string]string) bool {
	if len(r.labels) != len(labels) {
		return false
	}
	return findAllLabelsInLabels(r.labels, labels)
}
func findAllLabelsInLabels(labelsA map[string]string, labelsB map[string]string) bool {
	for k, lA := range labelsA {
		if lB, ok := labelsB[k]; !ok || lA != lB {
			return false
		}
	}
	return true
}

func (s *InMemoryFlatStorage) RemoveDefinition(featureName string) {
	s.Lock()
	defer s.Unlock()
	delete(s.table, featureName)
}

func (s *InMemoryFlatStorage) RemoveOverride(featureName string, labels map[string]string) {
	s.Lock()
	defer s.Unlock()
	if d, ok := s.table[featureName]; ok {
		d.removeRecords(labels)
		return
	}
	logrus.WithField("FeatureName", featureName).WithField("Action", "RemoveOverride").Error("Definition not found")
}
