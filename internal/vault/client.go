package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api     *vaultapi.Client
	Address string
	Token   string
}

// NewClient creates a new Vault client using environment variables or explicit config.
func NewClient(address, token string) (*Client, error) {
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		address = "http://127.0.0.1:8200"
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required: set VAULT_TOKEN or pass --token")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault api client: %w", err)
	}
	api.SetToken(token)

	return &Client{
		api:     api,
		Address: address,
		Token:   token,
	}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// path should be the logical path, e.g. "secret/data/myapp/config".
func (c *Client) ReadSecretVersion(path string, version int) (map[string]interface{}, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := c.api.Logical().ReadWithData(path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q (version %d): %w", path, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at %q (version %d)", path, version)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret data format at %q", path)
	}
	return data, nil
}
