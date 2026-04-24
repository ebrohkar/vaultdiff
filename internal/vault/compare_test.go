package vault

import (
	"testing"
)

func TestToStringMap_BasicTypes(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	result := toStringMap(input)

	if result["key1"] != "value1" {
		t.Errorf("expected key1=value1, got %s", result["key1"])
	}
	if result["key2"] != "42" {
		t.Errorf("expected key2=42, got %s", result["key2"])
	}
	if result["key3"] != "true" {
		t.Errorf("expected key3=true, got %s", result["key3"])
	}
}

func TestToStringMap_Empty(t *testing.T) {
	result := toStringMap(map[string]interface{}{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestVersionPair_ToStringMaps(t *testing.T) {
	vp := &VersionPair{
		Path: "secret/data/myapp",
		EnvA: "https://vault-dev.example.com",
		EnvB: "https://vault-prod.example.com",
		DataA: map[string]interface{}{
			"DB_HOST": "dev-db.internal",
			"DB_PORT": "5432",
		},
		DataB: map[string]interface{}{
			"DB_HOST": "prod-db.internal",
			"DB_PORT": "5432",
		},
	}

	a, b := vp.ToStringMaps()

	if a["DB_HOST"] != "dev-db.internal" {
		t.Errorf("expected dev-db.internal, got %s", a["DB_HOST"])
	}
	if b["DB_HOST"] != "prod-db.internal" {
		t.Errorf("expected prod-db.internal, got %s", b["DB_HOST"])
	}
	if a["DB_PORT"] != "5432" || b["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432 in both maps")
	}
}

func TestVersionPair_ToStringMaps_NilData(t *testing.T) {
	vp := &VersionPair{
		Path:  "secret/data/empty",
		DataA: nil,
		DataB: map[string]interface{}{"KEY": "val"},
	}

	a, b := vp.ToStringMaps()

	if len(a) != 0 {
		t.Errorf("expected empty map for nil DataA, got %d entries", len(a))
	}
	if b["KEY"] != "val" {
		t.Errorf("expected KEY=val in DataB map")
	}
}
