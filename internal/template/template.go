// Package template provides snapshot templating: apply variable substitution
// across a snapshot using a set of placeholder values.
package template

import (
	"fmt"
	"strings"
)

// Manager defines the interface for loading and saving snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the output of a template application.
type Result struct {
	Name    string
	Applied int
	Vars    map[string]string
}

// Run loads the snapshot identified by src, replaces all occurrences of each
// key in placeholders with its value inside every env var value, then saves
// the result as dst.
func Run(m Manager, src, dst string, placeholders map[string]string) (Result, error) {
	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("template: load %q: %w", src, err)
	}

	out := make(map[string]string, len(vars))
	applied := 0
	for k, v := range vars {
		newVal := v
		for placeholder, replacement := range placeholders {
			token := "{{" + placeholder + "}}"
			if strings.Contains(newVal, token) {
				newVal = strings.ReplaceAll(newVal, token, replacement)
				applied++
			}
		}
		out[k] = newVal
	}

	if err := m.Save(dst, out); err != nil {
		return Result{}, fmt.Errorf("template: save %q: %w", dst, err)
	}

	return Result{Name: dst, Applied: applied, Vars: out}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Saved snapshot %q with %d substitution(s)\n", r.Name, r.Applied)
	return sb.String()
}
