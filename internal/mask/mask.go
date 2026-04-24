// Package mask provides utilities for partially obscuring sensitive
// environment variable values in output, leaving a visible prefix so
// users can identify which credential is set without exposing it.
package mask

import "strings"

// DefaultVisibleChars is the number of leading characters shown before
// the mask is applied.
const DefaultVisibleChars = 4

// Options controls how masking is applied.
type Options struct {
	// VisibleChars is how many leading characters remain unmasked.
	VisibleChars int
	// MaskChar is the character used to fill the masked portion.
	MaskChar rune
	// MaskLen is the fixed length of the mask suffix (0 = use actual length).
	MaskLen int
}

// DefaultOptions returns sensible defaults for masking.
func DefaultOptions() Options {
	return Options{
		VisibleChars: DefaultVisibleChars,
		MaskChar:     '*',
		MaskLen:      8,
	}
}

// Value masks a single string value according to opts.
// If the value is shorter than or equal to VisibleChars the entire
// value is replaced by the mask so that length is not leaked.
func Value(v string, opts Options) string {
	if opts.MaskChar == 0 {
		opts.MaskChar = '*'
	}
	maskLen := opts.MaskLen
	if maskLen <= 0 {
		maskLen = len(v)
		if maskLen < opts.VisibleChars {
			maskLen = opts.VisibleChars
		}
	}
	if len(v) <= opts.VisibleChars {
		return strings.Repeat(string(opts.MaskChar), maskLen)
	}
	return v[:opts.VisibleChars] + strings.Repeat(string(opts.MaskChar), maskLen)
}

// Apply masks each value in vars, returning a new map.
// Keys are never modified.
func Apply(vars map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = Value(v, opts)
	}
	return out
}

// ApplyKeys masks only the values whose keys are listed in keys.
// All other entries are copied unmodified.
func ApplyKeys(vars map[string]string, keys []string, opts Options) map[string]string {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if _, ok := set[k]; ok {
			out[k] = Value(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}
