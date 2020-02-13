package test

import (
	"fmt"
	"testing"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage/tree"
)

type (
	storageConstructor func() storage.Storage
	action             func(s storage.Storage, t *testing.T)
	scenario           struct {
		name    string
		actions []action
	}
)

var currentStorageConstructor storageConstructor

func TestMain(m *testing.M) {
	var constructors = []storageConstructor{
		tree.NewInMemoryStorage,
	}
	for _, c := range constructors {
		currentStorageConstructor = c
		m.Run()
	}
}

func define(fName, fValue string) action {
	return func(s storage.Storage, t *testing.T) {
		s.Define(crd.FeatureFlag{FeatureID: fName, Value: fValue})
	}
}

func rmDef(key string) action {
	return func(s storage.Storage, t *testing.T) {
		s.RemoveDefinition(key)
	}
}

func rmOverride(key string, fields map[string]string) action {
	return func(s storage.Storage, t *testing.T) {
		s.RemoveOverride(key, fields)
	}
}

func override(fName, fValue, fOrigin string, priority int, labels map[string]string) action {
	return func(s storage.Storage, t *testing.T) {
		s.Override(crd.FeatureFlagOverride{FeatureID: fName, Value: fValue, OverrideName: fOrigin, Priority: priority, Labels: labels})
	}
}

func findAssert(key string, fields map[string]string, r *pb.FeatureFlag) action {
	return func(s storage.Storage, t *testing.T) {
		foundFF := s.Find(key, fields)
		if !isEquals(foundFF, r) {
			t.Errorf("Expected %#v but was %#v", r, foundFF)
		}
	}
}

func isEquals(ff1 *pb.FeatureFlag, ff2 *pb.FeatureFlag) bool {
	if ff1 == nil {
		if ff2 == nil {
			return true
		}
		return false
	}
	if ff2 == nil {
		return false
	}
	return ff1.Value == ff2.Value && ff1.FeatureName == ff2.FeatureName && ff1.Origin == ff2.Origin
}

func findAllAssert(fields map[string]string, ffs []*pb.FeatureFlag) action {
	return func(s storage.Storage, t *testing.T) {
		actualFFs := s.FindAll(fields)

		if len(actualFFs) != len(ffs) {
			t.Errorf("Expected len %+v but was %+v", len(ffs), len(actualFFs))
			return
		}
	expected:
		for _, ff := range ffs {
			found := false
			for _, actualFF := range actualFFs {
				if isEquals(ff, actualFF) {
					found = true
				}
			}
			if found {
				continue expected
			}
			t.Errorf("Expected %#v but not found", ff)
			return
		}
	}
}

func result(name, val, origin string) *pb.FeatureFlag {
	return &pb.FeatureFlag{FeatureName: name, Value: val, Origin: origin}
}

func defaultResult(name, val string) *pb.FeatureFlag {
	return &pb.FeatureFlag{FeatureName: name, Value: val}
}

func TestScenarios(t *testing.T) {
	for _, s := range scenarios {
		storageObj := currentStorageConstructor()
		t.Run(fmt.Sprintf("%s - %T", s.name, storageObj), func(t *testing.T) {
			for _, a := range s.actions {
				a(storageObj, t)
			}
		})
	}
}

var scenarios = []scenario{
	{
		name: "store definition - found default",
		actions: []action{
			define("key1", "value1"),
			findAssert("key1", nil, defaultResult("key1", "value1")),
		},
	}, {
		name: "store definition twice - found second definition",
		actions: []action{
			define("key1", "value1"),
			define("key1", "value2"),
			findAssert("key1", nil, defaultResult("key1", "value2")),
		},
	}, {
		name: "store definition and override - found override",
		actions: []action{
			define("key1", "value1"),
			override("key1", "overrideValue1", "origin1", 1, map[string]string{
				"field1": "fieldValue1",
			}),
			findAssert("key1", map[string]string{
				"field1": "fieldValue1",
			}, result("key1", "overrideValue1", "origin1")),
		},
	}, {
		name: "store definition - unknown not found",
		actions: []action{
			define("key1", "value1"),
			findAssert("key_unknown", nil, nil),
		},
	}, {
		name: "store definition and override multiple labels - found override 2",
		actions: []action{
			define("key1", "value1"),
			override("key1", "overrideValue1", "origin1", 1, map[string]string{
				"field1": "fieldValue1",
			}),
			override("key1", "overrideValue2", "origin2", 2, map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}),
			findAssert("key1", map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}, result("key1", "overrideValue2", "origin2")),
		},
	}, {
		name: "store definition and override multiple labels with less priority - found override 1",
		actions: []action{
			define("key1", "value1"),
			override("key1", "overrideValue1", "origin1", 2, map[string]string{
				"field1": "fieldValue1",
			}),
			override("key1", "overrideValue2", "origin2", 1, map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}),
			findAssert("key1", map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}, result("key1", "overrideValue1", "origin1")),
		},
	}, {
		name: "store definition and override and remove definition - not found",
		actions: []action{
			define("key1", "value1"),
			override("key1", "overrideValue1", "origin1", 1, map[string]string{
				"field1": "fieldValue1",
			}),
			rmDef("key1"),
			findAssert("key1", map[string]string{
				"field1": "fieldValue1",
			}, nil),
			findAssert("key1", nil, nil),
		},
	}, {
		name: "store definition and override and remove override - not found",
		actions: []action{
			define("key1", "value1"),
			override("key1", "overrideValue1", "origin1", 1, map[string]string{
				"field1": "fieldValue1",
			}),
			rmOverride("key1", map[string]string{
				"field1": "fieldValue1",
			}),
			findAssert("key1", nil, defaultResult("key1", "value1")),
			findAssert("key1", map[string]string{
				"field1": "fieldValue1",
			}, defaultResult("key1", "value1")),
		},
	}, {
		name: "store definition and override multiple labels - found all",
		actions: []action{
			define("key1", "value1"),
			define("key2", "value2_1"),
			override("key1", "overrideValue1", "origin1", 1, map[string]string{
				"field1": "fieldValue1",
			}),
			override("key1", "overrideValue2", "origin2", 100, map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}),
			override("key2", "overrideValue2_1", "origin1", 100, map[string]string{
				"field1": "fieldValue1",
			}),
			override("key2", "overrideValue2_2", "origin2", 1, map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}),
			findAllAssert(nil, []*pb.FeatureFlag{
				defaultResult("key1", "value1"),
				defaultResult("key2", "value2_1"),
			}),
			findAllAssert(map[string]string{
				"field1": "fieldValue1",
			}, []*pb.FeatureFlag{
				result("key1", "overrideValue1", "origin1"),
				result("key2", "overrideValue2_1", "origin1"),
			}),
			findAllAssert(map[string]string{
				"field1": "fieldValue1",
				"field2": "fieldValue2",
			}, []*pb.FeatureFlag{
				result("key1", "overrideValue2", "origin2"),
				result("key2", "overrideValue2_1", "origin1"),
			}),
		},
	},
}
