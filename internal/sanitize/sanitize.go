// Package sanitize provides utilities for cleaning and normalising
// environment variable maps before they are stored as snapshots.
package sanitize

import (
	"strings"
)

// Options controls which sanitisation passes are applied.
type Options struct {
	// TrimSpace removes leading and trailing whitespace from values.
	TrimSpace bool
	// RemoveEmpty drops keys whose value is the empty string (after trimming).
	RemoveEmpty bool
	// UpperKeys normalises every key to upper-case.
	UpperKeys bool
	// StripPrefix removes a common prefix from every key, if set.
	StripPrefix string
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:   true,
		RemoveEmpty: false,
		UpperKeys:   false,
	}
}

// Run applies the requested sanitisation passes to a copy of vars and
// returns the cleaned map together with a summary of changes made.
func Run(vars map[string]string, opts Options) (map[string]string, Result) {
	out := make(map[string]string, len(vars))
	var res Result

	for k, v := range vars {
		origKey := k
		origVal := v

		if opts.UpperKeys {
			k = strings.ToUpper(k)
			if k != origKey {
				res.RenamedKeys = append(res.RenamedKeys, origKey)
			}
		}

		if opts.StripPrefix != "" && strings.HasPrefix(k, opts.StripPrefix) {
			k = k[len(opts.StripPrefix):]
			if k == "" {
				res.DroppedKeys = append(res.DroppedKeys, origKey)
				continue
			}
		}

		if opts.TrimSpace {
			v = strings.TrimSpace(v)
			if v != origVal {
				res.TrimmedKeys = append(res.TrimmedKeys, origKey)
			}
		}

		if opts.RemoveEmpty && v == "" {
			res.DroppedKeys = append(res.DroppedKeys, origKey)
			continue
		}

		out[k] = v
	}

	return out, res
}

// Result summarises what sanitisation did to the input map.
type Result struct {
	TrimmedKeys []string
	RenamedKeys []string
	DroppedKeys []string
}

// Changed reports whether any modification was made.
func (r Result) Changed() bool {
	return len(r.TrimmedKeys) > 0 || len(r.RenamedKeys) > 0 || len(r.DroppedKeys) > 0
}
