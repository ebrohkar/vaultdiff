package vault

import (
	"context"
	"fmt"
	"sort"
)

// VersionMeta holds metadata about a specific secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListSecretVersions returns metadata for all versions of a KV v2 secret.
func (c *Client) ListSecretVersions(ctx context.Context, mount, path string) ([]VersionMeta, error) {
	secretPath := fmt.Sprintf("%s/metadata/%s", mount, path)

	secret, err := c.logical.ReadWithContext(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("listing versions for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for secret %q", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata for %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for %q", path)
	}

	var metas []VersionMeta
	for versionStr, vRaw := range versionsMap {
		vMap, ok := vRaw.(map[string]interface{})
		if !ok {
			continue
		}

		var num int
		fmt.Sscanf(versionStr, "%d", &num)

		meta := VersionMeta{
			Version:     num,
			Destroyed:   boolVal(vMap, "destroyed"),
			CreatedTime: stringVal(vMap, "created_time"),
			DeletionTime: stringVal(vMap, "deletion_time"),
		}
		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	return metas, nil
}

func boolVal(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	b, _ := v.(bool)
	return b
}

func stringVal(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
