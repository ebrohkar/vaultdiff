package vault

import "fmt"

// MergeOptions configures the merge operation.
type MergeOptions struct {
	Mount      string
	SourcePath string
	DestPath   string
	// Strategy defines how conflicting keys are resolved.
	// "source" overwrites dest keys with source values (default).
	// "dest" keeps dest values when keys conflict.
	// "error" aborts if any key conflict is detected.
	Strategy string
}

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	SourcePath  string
	DestPath    string
	Merged      map[string]string
	Conflicts   []string
	AddedKeys   []string
	SkippedKeys []string
}

// DefaultMergeOptions returns MergeOptions with sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Mount:    "secret",
		Strategy: "source",
	}
}

// MergeSecret merges key/value pairs from sourcePath into destPath using the
// provided strategy. It does not write back to Vault; callers are responsible
// for persisting the returned MergeResult.Merged map.
func MergeSecret(client *Client, opts MergeOptions) (*MergeResult, error) {
	if opts.SourcePath == "" {
		return nil, fmt.Errorf("merge: source path must not be empty")
	}
	if opts.DestPath == "" {
		return nil, fmt.Errorf("merge: dest path must not be empty")
	}
	if opts.SourcePath == opts.DestPath {
		return nil, fmt.Errorf("merge: source and dest paths must differ")
	}
	if client == nil {
		return nil, fmt.Errorf("merge: client must not be nil")
	}

	validStrategies := map[string]bool{"source": true, "dest": true, "error": true}
	if !validStrategies[opts.Strategy] {
		return nil, fmt.Errorf("merge: unknown strategy %q; must be source, dest, or error", opts.Strategy)
	}

	srcData, err := client.ReadSecretVersion(opts.Mount, opts.SourcePath, 0)
	if err != nil {
		return nil, fmt.Errorf("merge: reading source: %w", err)
	}
	dstData, err := client.ReadSecretVersion(opts.Mount, opts.DestPath, 0)
	if err != nil {
		return nil, fmt.Errorf("merge: reading dest: %w", err)
	}

	src := toStringMap(srcData)
	dst := toStringMap(dstData)

	result := &MergeResult{
		SourcePath: opts.SourcePath,
		DestPath:   opts.DestPath,
		Merged:     make(map[string]string),
	}

	for k, v := range dst {
		result.Merged[k] = v
	}

	for k, v := range src {
		if existing, exists := result.Merged[k]; exists {
			result.Conflicts = append(result.Conflicts, k)
			switch opts.Strategy {
			case "error":
				return nil, fmt.Errorf("merge: conflict on key %q (dest=%q, source=%q)", k, existing, v)
			case "dest":
				result.SkippedKeys = append(result.SkippedKeys, k)
				continue
			default: // "source"
				result.Merged[k] = v
			}
		} else {
			result.Merged[k] = v
			result.AddedKeys = append(result.AddedKeys, k)
		}
	}

	return result, nil
}
