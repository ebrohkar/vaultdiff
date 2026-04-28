package vault

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshalEnv_MultipleKeys(t *testing.T) {
	data := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	out := string(marshalEnv(data))
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("missing DB_HOST in env output: %q", out)
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("missing DB_PORT in env output: %q", out)
	}
}

func TestMarshalEnv_EmptyMap(t *testing.T) {
	out := marshalEnv(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output for empty map, got %q", out)
	}
}

func TestMarshalEnv_SpecialCharacters(t *testing.T) {
	data := map[string]string{"KEY": "val=with=equals"}
	out := string(marshalEnv(data))
	if !strings.Contains(out, "KEY=val=with=equals") {
		t.Errorf("unexpected env output: %q", out)
	}
}

func TestExportFormatJSON_ValidJSON(t *testing.T) {
	data := map[string]string{"username": "admin", "password": "s3cr3t"}
	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded["username"] != "admin" {
		t.Errorf("unexpected username: %s", decoded["username"])
	}
}

func TestExportFormat_Constants(t *testing.T) {
	if ExportFormatJSON != "json" {
		t.Errorf("ExportFormatJSON should be 'json', got %q", ExportFormatJSON)
	}
	if ExportFormatEnv != "env" {
		t.Errorf("ExportFormatEnv should be 'env', got %q", ExportFormatEnv)
	}
}
