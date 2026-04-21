package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_WithToken(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Token != "test-token" {
		t.Errorf("expected token %q, got %q", "test-token", client.Token)
	}
	if client.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected address %q, got %q", "http://127.0.0.1:8200", client.Address)
	}
}

func TestReadSecretVersion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"data": {"username": "admin", "password": "s3cr3t"},
				"metadata": {"version": 1}
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "fake-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	data, err := client.ReadSecretVersion("secret/data/myapp", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["username"] != "admin" {
		t.Errorf("expected username %q, got %q", "admin", data["username"])
	}
}

func TestReadSecretVersion_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "fake-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.ReadSecretVersion("secret/data/missing", 0)
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
