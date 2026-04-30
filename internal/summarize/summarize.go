// Package summarize provides a high-level summary of a snapshot's contents.
package summarize

import (
	"fmt"
	"sort"
	"strings"
)

// Manager is the interface required to load a snapshot's variables.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Result holds the computed summary for a snapshot.
type Result struct {
	Name       string
	Total      int
	Empty      int
	Sensitive  int
	Prefixes   map[string]int // top-level prefix (before first '_') -> count
}

// sensitivePatterns are substrings that suggest a key is sensitive.
var sensitivePatterns = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "CREDENTIAL"}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Run loads the named snapshot and computes a summary.
func Run(m Manager, name string) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("summarize: %w", err)
	}

	r := Result{
		Name:     name,
		Total:    len(vars),
		Prefixes: make(map[string]int),
	}

	for k, v := range vars {
		if v == "" {
			r.Empty++
		}
		if isSensitive(k) {
			r.Sensitive++
		}
		parts := strings.SplitN(k, "_", 2)
		r.Prefixes[parts[0]]++
	}

	return r, nil
}

// Format renders a Result as a human-readable string.
func Format(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot : %s\n", r.Name)
	fmt.Fprintf(&sb, "Total    : %d\n", r.Total)
	fmt.Fprintf(&sb, "Empty    : %d\n", r.Empty)
	fmt.Fprintf(&sb, "Sensitive: %d\n", r.Sensitive)

	if len(r.Prefixes) > 0 {
		sb.WriteString("Prefixes :\n")
		keys := make([]string, 0, len(r.Prefixes))
		for k := range r.Prefixes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %-20s %d\n", k, r.Prefixes[k])
		}
	}
	return sb.String()
}
