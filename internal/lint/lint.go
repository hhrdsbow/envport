// Package lint checks snapshots for common issues such as empty values,
// keys with suspicious characters, or duplicate entries.
package lint

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single lint finding.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%-24s %s", i.Key, i.Message)
}

// Manager is the interface lint needs to load a snapshot's vars.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Run loads the named snapshot and returns any lint issues found.
func Run(m Manager, name string) ([]Issue, error) {
	vars, err := m.Load(name)
	if err != nil {
		return nil, fmt.Errorf("lint: load %q: %w", name, err)
	}

	var issues []Issue
	seen := make(map[string]bool)

	for k, v := range vars {
		if seen[k] {
			issues = append(issues, Issue{Key: k, Message: "duplicate key"})
		}
		seen[k] = true

		if strings.TrimSpace(k) == "" {
			issues = append(issues, Issue{Key: k, Message: "blank key name"})
			continue
		}

		if strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Key: k, Message: "empty value"})
		}

		if err := validateKey(k); err != nil {
			issues = append(issues, Issue{Key: k, Message: err.Error()})
		}
	}

	return issues, nil
}

// Format returns a human-readable summary of the issues slice.
func Format(issues []Issue) string {
	if len(issues) == 0 {
		return "no issues found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString(iss.String())
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func validateKey(k string) error {
	for _, r := range k {
		if unicode.IsSpace(r) {
			return fmt.Errorf("key contains whitespace")
		}
		if r == '=' {
			return fmt.Errorf("key contains '='")
		}
	}
	return nil
}
