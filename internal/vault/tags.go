package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretTag represents a key-value metadata tag attached to a secret path.
type SecretTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TagResult holds the outcome of a tag operation.
type TagResult struct {
	Path    string      `json:"path"`
	Tags    []SecretTag `json:"tags"`
	Updated bool        `json:"updated"`
}

// GetTags retrieves custom metadata tags for a KV v2 secret path.
func GetTags(ctx context.Context, client *Client, mount, path string) ([]SecretTag, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("secret path must not be empty")
	}
	if strings.TrimSpace(mount) == "" {
		mount = "secret"
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)
	secret, err := client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return []SecretTag{}, nil
	}

	customMeta, ok := secret.Data["custom_metadata"].(map[string]interface{})
	if !ok {
		return []SecretTag{}, nil
	}

	tags := make([]SecretTag, 0, len(customMeta))
	for k, v := range customMeta {
		tags = append(tags, SecretTag{Key: k, Value: fmt.Sprintf("%v", v)})
	}
	return tags, nil
}

// SetTags writes custom metadata tags to a KV v2 secret path.
func SetTags(ctx context.Context, client *Client, mount, path string, tags []SecretTag) (*TagResult, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("secret path must not be empty")
	}
	if strings.TrimSpace(mount) == "" {
		mount = "secret"
	}

	customMeta := make(map[string]interface{}, len(tags))
	for _, t := range tags {
		if strings.TrimSpace(t.Key) == "" {
			return nil, fmt.Errorf("tag key must not be empty")
		}
		customMeta[t.Key] = t.Value
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)
	_, err := client.Logical().WriteWithContext(ctx, metaPath, map[string]interface{}{
		"custom_metadata": customMeta,
	})
	if err != nil {
		return nil, fmt.Errorf("writing tags for %q: %w", path, err)
	}

	return &TagResult{Path: path, Tags: tags, Updated: true}, nil
}
