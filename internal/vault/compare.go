package vault

import (
	"context"
	"fmt"
)

// VersionPair holds two versions of a secret for comparison.
type VersionPair struct {
	Path    string
	EnvA    string
	EnvB    string
	DataA   map[string]interface{}
	DataB   map[string]interface{}
}

// FetchVersionPair retrieves the same secret path from two different clients
// (representing two environments or two version numbers) and returns a VersionPair.
func FetchVersionPair(
	ctx context.Context,
	clientA *Client,
	clientB *Client,
	path string,
	versionA int,
	versionB int,
) (*VersionPair, error) {
	dataA, err := clientA.ReadSecretVersion(ctx, path, versionA)
	if err != nil {
		return nil, fmt.Errorf("fetching version %d from env %s: %w", versionA, clientA.Address(), err)
	}

	dataB, err := clientB.ReadSecretVersion(ctx, path, versionB)
	if err != nil {
		return nil, fmt.Errorf("fetching version %d from env %s: %w", versionB, clientB.Address(), err)
	}

	return &VersionPair{
		Path:  path,
		EnvA:  clientA.Address(),
		EnvB:  clientB.Address(),
		DataA: dataA,
		DataB: dataB,
	}, nil
}

// ToStringMaps converts the VersionPair data fields to map[string]string
// for use with the diff.Compare function.
func (vp *VersionPair) ToStringMaps() (map[string]string, map[string]string) {
	a := toStringMap(vp.DataA)
	b := toStringMap(vp.DataB)
	return a, b
}

func toStringMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}
