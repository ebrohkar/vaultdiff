package vault

import (
	"context"
	"testing"
)

func TestSecretTag_Fields(t *testing.T) {
	tag := SecretTag{Key: "env", Value: "production"}
	if tag.Key != "env" {
		t.Errorf("expected key %q, got %q", "env", tag.Key)
	}
	if tag.Value != "production" {
		t.Errorf("expected value %q, got %q", "production", tag.Value)
	}
}

func TestTagResult_Fields(t *testing.T) {
	result := &TagResult{
		Path:    "myapp/config",
		Tags:    []SecretTag{{Key: "team", Value: "platform"}},
		Updated: true,
	}
	if result.Path != "myapp/config" {
		t.Errorf("expected path %q, got %q", "myapp/config", result.Path)
	}
	if !result.Updated {
		t.Error("expected Updated to be true")
	}
	if len(result.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(result.Tags))
	}
	if result.Tags[0].Key != "team" {
		t.Errorf("expected tag key %q, got %q", "team", result.Tags[0].Key)
	}
}

func TestGetTags_EmptyPath(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	_, err = GetTags(context.Background(), client, "secret", "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestSetTags_EmptyPath(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	_, err = SetTags(context.Background(), client, "secret", "", []SecretTag{{Key: "env", Value: "dev"}})
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestSetTags_EmptyTagKey(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	_, err = SetTags(context.Background(), client, "secret", "myapp/config", []SecretTag{{Key: "", Value: "dev"}})
	if err == nil {
		t.Fatal("expected error for empty tag key, got nil")
	}
}

func TestSetTags_DefaultMount(t *testing.T) {
	// Validates that an empty mount falls back to "secret" without panicking.
	// A real Vault call would fail; we only check the validation path here.
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	// We expect a network error (not a validation error) since mount defaults.
	_, err = SetTags(context.Background(), client, "", "myapp/config", []SecretTag{{Key: "env", Value: "dev"}})
	if err == nil {
		t.Fatal("expected network/write error, got nil")
	}
	// Ensure it is NOT a validation error about empty path or key.
	if err.Error() == "secret path must not be empty" || err.Error() == "tag key must not be empty" {
		t.Errorf("unexpected validation error: %v", err)
	}
}
