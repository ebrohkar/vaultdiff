package vault

import (
	"fmt"
	"strings"
)

// SearchOptions configures a secret search operation.
type SearchOptions struct {
	Mount      string
	PathPrefix string
	Keyword    string
	CaseSensitive bool
}

// SearchResult holds a matched secret path and the matching keys within it.
type SearchResult struct {
	Path        string
	MatchedKeys []string
}

// SearchSecrets walks secrets under PathPrefix and returns paths whose data
// contains keys or values matching Keyword.
func SearchSecrets(client *Client, opts SearchOptions) ([]SearchResult, error) {
	if opts.Mount == "" {
		opts.Mount = "secret"
	}
	if opts.Keyword == "" {
		return nil, fmt.Errorf("search keyword must not be empty")
	}

	paths, err := listAllPaths(client, opts.Mount, opts.PathPrefix)
	if err != nil {
		return nil, fmt.Errorf("listing paths: %w", err)
	}

	var results []SearchResult
	for _, p := range paths {
		data, err := client.ReadSecretVersion(p, 0)
		if err != nil || data == nil {
			continue
		}
		matched := matchingKeys(data, opts.Keyword, opts.CaseSensitive)
		if len(matched) > 0 {
			results = append(results, SearchResult{Path: p, MatchedKeys: matched})
		}
	}
	return results, nil
}

// listAllPaths recursively lists KV paths under mount/prefix.
func listAllPaths(client *Client, mount, prefix string) ([]string, error) {
	listPath := mount + "/metadata/" + prefix
	secret, err := client.Logical().List(listPath)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}
	keys, _ := secret.Data["keys"].([]interface{})
	var paths []string
	for _, k := range keys {
		key, _ := k.(string)
		if strings.HasSuffix(key, "/") {
			sub, err := listAllPaths(client, mount, prefix+key)
			if err == nil {
				paths = append(paths, sub...)
			}
		} else {
			paths = append(paths, mount+"/data/"+prefix+key)
		}
	}
	return paths, nil
}

// matchingKeys returns keys (and keys whose values) contain the keyword.
func matchingKeys(data map[string]interface{}, keyword string, caseSensitive bool) []string {
	var matched []string
	needle := keyword
	if !caseSensitive {
		needle = strings.ToLower(keyword)
	}
	for k, v := range data {
		haystack := k + fmt.Sprintf("%v", v)
		if !caseSensitive {
			haystack = strings.ToLower(haystack)
		}
		if strings.Contains(haystack, needle) {
			matched = append(matched, k)
		}
	}
	return matched
}
