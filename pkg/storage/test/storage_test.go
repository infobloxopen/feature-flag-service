package test

import (
	"fmt"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage/flat"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage/tree"
	"testing"
)

type implementation struct {
	name        string
	constructor func() storage.Storage
}

var currentImplementation *implementation

func TestMain(m *testing.M) {
	var implementations = []implementation{
		{"Flat", flat.NewInMemoryStorage},
		{"Tree", tree.NewInMemoryStorage},
	}
	for _, i := range implementations {
		currentImplementation = &i
		m.Run()
	}
}

func assertEquals(t *testing.T, actualR, expectedR *pb.FeatureFlag) {
	if expectedR == nil && actualR == nil {
		return
	} else if expectedR == nil && actualR != nil {
		t.Error("Expected nil but was not nil")
	} else if expectedR != nil && actualR == nil {
		t.Error("Expected value not nil but was nil")
	} else if actualR.FeatureName != expectedR.FeatureName {
		t.Errorf("Expected FeatureName %q but was %q", expectedR.FeatureName, actualR.FeatureName)
	} else if actualR.Value != expectedR.Value {
		t.Errorf("Expected value %q but was %q", expectedR.Value, actualR.Value)
	} else if actualR.Origin != expectedR.Origin {
		t.Errorf("Expected Origin %q but was %q", expectedR.Origin, actualR.Origin)
	}
}

type (
	TestFeatureSearch struct {
		id     string
		labels map[string]string
	}
	TestCase struct {
		name           string
		definitions    []storage.FeatureFlagDefinition
		overrides      []storage.FeatureFlagOverride
		search         TestFeatureSearch
		expectedValues *pb.FeatureFlag
	}
)

var testTable = []TestCase{
	{
		name: "store definition - found default",
		definitions: []storage.FeatureFlagDefinition{
			{FeatureName: "feature", DefaultValue: "value"},
		},
		search:         TestFeatureSearch{id: "feature"},
		expectedValues: &pb.FeatureFlag{FeatureName: "feature", Value: "value", Origin: ""},
	}, {
		name: "store definition twice - found default value with value from second definition",
		definitions: []storage.FeatureFlagDefinition{
			{FeatureName: "feature", DefaultValue: "value"},
			{FeatureName: "feature", DefaultValue: "value2"},
		},
		search:         TestFeatureSearch{id: "feature"},
		expectedValues: &pb.FeatureFlag{FeatureName: "feature", Value: "value2", Origin: ""},
	}, {
		name: "store definition - unknown not found",
		definitions: []storage.FeatureFlagDefinition{
			{FeatureName: "feature", DefaultValue: "value"},
		},
		search:         TestFeatureSearch{id: "unknown_feature"},
		expectedValues: nil,
	}, {
		name: "store definition and override - retrieve with labels found override",
		definitions: []storage.FeatureFlagDefinition{
			{FeatureName: "feature", DefaultValue: "value"},
		},
		overrides: []storage.FeatureFlagOverride{
			{FeatureName: "feature", Origin: "override_name", Priority: 42, Value: "override value", Labels: map[string]string{
				"label":  "label_value",
				"label2": "label_value2",
			}},
		},
		search: TestFeatureSearch{id: "feature", labels: map[string]string{
			"label":  "label_value",
			"label2": "label_value2",
		}},
		expectedValues: &pb.FeatureFlag{FeatureName: "feature", Value: "override value", Origin: "override_name"},
	},
}

func TestStorage_SaveAndFind(t *testing.T) {
	for _, testCase := range testTable {
		t.Run(fmt.Sprintf("%s - %s", testCase.name, currentImplementation.name), func(t *testing.T) {
			t.Log(testCase.name)
			storage := currentImplementation.constructor()
			for _, definition := range testCase.definitions {
				storage.Define(definition)
			}
			for _, override := range testCase.overrides {
				storage.Override(override)
			}
			result := storage.Find(testCase.search.id, testCase.search.labels)
			assertEquals(t, result, testCase.expectedValues)
		})
	}
}

func TestInMemoryTreeStorage_RemoveDefinition(t *testing.T) {
	s := currentImplementation.constructor()
	s.Define(storage.FeatureFlagDefinition{FeatureName: "feature", DefaultValue: "value"})
	s.Override(storage.FeatureFlagOverride{FeatureName: "feature", Origin: "override_name", Priority: 42, Value: "override value", Labels: map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	}})
	if ff := s.Find("feature", nil); ff == nil {
		t.Error("Expected not nil")
	} else if ff = s.Find("feature", map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	}); ff == nil {
		t.Error("Expected not nil")
	}

	s.RemoveDefinition("feature")

	if ff := s.Find("feature", nil); ff != nil {
		t.Error("Expected not nil")
	} else if ff = s.Find("feature", map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	}); ff != nil {
		t.Error("Expected not nil")
	}
}

func TestInMemoryTreeStorage_RemoveOverride(t *testing.T) {
	s := currentImplementation.constructor()
	s.Define(storage.FeatureFlagDefinition{FeatureName: "feature", DefaultValue: "value"})
	s.Override(storage.FeatureFlagOverride{FeatureName: "feature", Origin: "override_name", Priority: 42, Value: "override value", Labels: map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	}})
	if ff := s.Find("feature", nil); ff == nil {
		t.Error("Expected not nil")
	} else if ff = s.Find("feature", map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	}); ff == nil {
		t.Error("Expected not nil")
	}

	s.RemoveOverride("feature", map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	})

	if ffDef := s.Find("feature", nil); ffDef == nil {
		t.Error("Expected not nil")
	} else if ffOver := s.Find("feature", map[string]string{
		"label":  "label_value",
		"label2": "label_value2",
	}); ffOver == nil {
		t.Error("Expected nil")
	} else {
		assertEquals(t, ffOver, &pb.FeatureFlag{FeatureName: "feature", Value: "value"})
	}
}

func TestInMemoryTreeStorage_RemoveOverrideComplex(t *testing.T) {
	s := currentImplementation.constructor()
	s.Define(storage.FeatureFlagDefinition{FeatureName: "feature", DefaultValue: "value"})
	s.Override(storage.FeatureFlagOverride{FeatureName: "feature", Origin: "override_name1", Priority: 42, Value: "override value1", Labels: map[string]string{
		"account": "40",
	}})
	s.Override(storage.FeatureFlagOverride{FeatureName: "feature", Origin: "override_name2", Priority: 42, Value: "override value2", Labels: map[string]string{
		"account": "40",
		"user":    "5",
	}})

	s.RemoveOverride("feature", map[string]string{
		"account": "40",
	})

	ff := s.Find("feature", map[string]string{
		"account": "40",
		"user":    "5",
	})

	assertEquals(t, ff, &pb.FeatureFlag{FeatureName: "feature", Origin: "override_name2", Value: "override value2"})
}
