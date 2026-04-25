package vault

import (
	"context"
	"fmt"
	"strings"
)

// PolicyAccess represents the access level a token has to a secret path.
type PolicyAccess struct {
	Path        string
	CanRead     bool
	CanList     bool
	Denied      bool
}

// CheckPathAccess verifies whether the configured Vault token has read and list
// capabilities on the given secret path by calling the sys/capabilities-self endpoint.
func (c *Client) CheckPathAccess(ctx context.Context, path string) (*PolicyAccess, error) {
	body := map[string]interface{}{
		"paths": []string{path},
	}

	secret, err := c.logical.WriteWithContext(ctx, "sys/capabilities-self", body)
	if err != nil {
		return nil, fmt.Errorf("capabilities check failed for path %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("empty response from capabilities-self for path %q", path)
	}

	access := &PolicyAccess{Path: path}

	raw, ok := secret.Data[path]
	if !ok {
		access.Denied = true
		return access, nil
	}

	caps, ok := raw.([]interface{})
	if !ok {
		access.Denied = true
		return access, nil
	}

	for _, c := range caps {
		switch strings.ToLower(fmt.Sprintf("%v", c)) {
		case "read":
			access.CanRead = true
		case "list":
			access.CanList = true
		case "deny":
			access.Denied = true
		}
	}

	return access, nil
}
