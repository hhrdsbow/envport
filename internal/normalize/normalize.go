// Package normalize provides utilities for normalizing environment variable
// maps across snapshots — trimming whitespace, fixing key casing, and
// removing empty entries.
package normalize

import (
	"fmt"
	"strings"
)

// Options controls which normalization steps are applied.
type Options struct {
	UpperKeys   bool
	TrimValues  bool
	RemoveEmpty bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		UpperKeys:   true,
		TrimValues:  true,
		RemoveEmpty: false,
	}
}

// Result holds the output of a normalization run.
type Result struct {
	Vars    map[string]string
	Changed bool
}

// Run applies the given Options to vars and returns a Result.
// The original map is never mutated.
func Run(vars map[string]string, opts Options) Result {
	out := make(map[string]string, len(vars))
	changed := false

	for k, v := range vars {
		newKey := k
		if opts.UpperKeys {
			newKey = strings.ToUpper(k)
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
		}

		if opts.RemoveEmpty && newVal == "" {
			changed = true
			continue
		}

		if newKey != k || newVal != v {
			changed = true
		}

		out[newKey] = newVal
	}

	return Result{Vars: out, Changed: changed}
}

// Format returns a human-readable summary of normalization changes.
func Format(before, after map[string]string) string {
	var sb strings.Builder

	removed := 0
	for k := range before {
		if _, ok := after[k]; !ok {
			// only count as removed if the upper-cased key is also absent
			if _, ok2 := after[strings.ToUpper(k)]; !ok2 {
				removed++
			}
		}
	}
	if removed > 0 {
		sb.WriteString(fmt.Sprintf("removed %d empty key(s)\n", removed))
	}

	for k, v := range after {
		old, existed := before[k]
		if !existed {
			sb.WriteString(fmt.Sprintf("renamed key -> %s\n", k))
		} else if old != v {
			sb.WriteString(fmt.Sprintf("trimmed value for %s\n", k))
		}
	}

	return sb.String()
}
