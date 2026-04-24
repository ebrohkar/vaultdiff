package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

// mockLogical is a minimal stub for the Vault logical client used in tests.
type mockLogical struct {
	readFn func(path string) (*vaultapi.Secret, error)
}

func (m *mockLogical) ReadWithContext(_ context.Context, path string) (*vaultapi.Secret, error) {
	return m.readFn(path)
}

func TestListSecretVersions_Success(t *testing.T) {
	versionsData := map[string]interface{}{
		"1": map[string]interface{}{
			"created_time":  "2024-01-01T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
		"2": map[string]interface{}{
			"created_time":  "2024-02-01T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
	}

	client := &Client{
		logical: &mockLogical{
			readFn: func(_ string) (*vaultapi.Secret, error) {
				return &vaultapi.Secret{Data: map[string]interface{}{"versions": versionsData}}, nil
			},
		},
	}

	metas, err := client.ListSecretVersions(context.Background(), "secret", "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(metas))
	}
	if metas[0].Version != 1 || metas[1].Version != 2 {
		t.Errorf("versions not sorted: got %d, %d", metas[0].Version, metas[1].Version)
	}
	if metas[0].CreatedTime != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected created_time: %s", metas[0].CreatedTime)
	}
}

func TestListSecretVersions_NotFound(t *testing.T) {
	client := &Client{
		logical: &mockLogical{
			readFn: func(_ string) (*vaultapi.Secret, error) {
				return nil, nil
			},
		},
	}

	_, err := client.ListSecretVersions(context.Background(), "secret", "missing/path")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestListSecretVersions_Destroyed(t *testing.T) {
	versionsData := map[string]interface{}{
		"3": map[string]interface{}{
			"created_time":  "2024-03-01T00:00:00Z",
			"deletion_time": "2024-03-15T00:00:00Z",
			"destroyed":     true,
		},
	}

	client := &Client{
		logical: &mockLogical{
			readFn: func(_ string) (*vaultapi.Secret, error) {
				return &vaultapi.Secret{Data: map[string]interface{}{"versions": versionsData}}, nil
			},
		},
	}

	metas, err := client.ListSecretVersions(context.Background(), "secret", "myapp/old")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !metas[0].Destroyed {
		t.Error("expected version to be marked as destroyed")
	}
}
