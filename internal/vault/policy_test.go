package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestCheckPathAccess_ReadAndList(t *testing.T) {
	server := newMockVaultServer(t, map[string]mockResponse{
		"/v1/sys/capabilities-self": {
			statusCode: 200,
			body: map[string]interface{}{
				"secret/data/myapp": []interface{}{"read", "list"},
			},
		},
	})
	defer server.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = server.URL
	client, err := NewClient(cfg, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	access, err := client.CheckPathAccess(context.Background(), "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !access.CanRead {
		t.Error("expected CanRead to be true")
	}
	if !access.CanList {
		t.Error("expected CanList to be true")
	}
	if access.Denied {
		t.Error("expected Denied to be false")
	}
}

func TestCheckPathAccess_Denied(t *testing.T) {
	server := newMockVaultServer(t, map[string]mockResponse{
		"/v1/sys/capabilities-self": {
			statusCode: 200,
			body: map[string]interface{}{
				"secret/data/restricted": []interface{}{"deny"},
			},
		},
	})
	defer server.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = server.URL
	client, err := NewClient(cfg, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	access, err := client.CheckPathAccess(context.Background(), "secret/data/restricted")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !access.Denied {
		t.Error("expected Denied to be true")
	}
	if access.CanRead {
		t.Error("expected CanRead to be false")
	}
}

func TestCheckPathAccess_MissingPath(t *testing.T) {
	server := newMockVaultServer(t, map[string]mockResponse{
		"/v1/sys/capabilities-self": {
			statusCode: 200,
			body: map[string]interface{}{},
		},
	})
	defer server.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = server.URL
	client, err := NewClient(cfg, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	access, err := client.CheckPathAccess(context.Background(), "secret/data/unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !access.Denied {
		t.Error("expected Denied to be true when path not in response")
	}
}
