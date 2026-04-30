// Package typecheck validates that environment variable values conform to
// declared types such as int, float, bool, or url.
package typecheck

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Type represents a declared value type.
type Type string

const (
	TypeInt   Type = "int"
	TypeFloat Type = "float"
	TypeBool  Type = "bool"
	TypeURL   Type = "url"
)

// Violation describes a single type mismatch.
type Violation struct {
	Key      string
	Value    string
	Expected Type
	Reason   string
}

// Result holds the outcome of a type-check run.
type Result struct {
	Violations []Violation
}

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Run checks the vars in the named snapshot against the provided type map.
// typeMap is key -> Type; keys not present in typeMap are skipped.
func Run(m Manager, name string, typeMap map[string]Type) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("typecheck: load %q: %w", name, err)
	}

	var violations []Violation
	for key, expected := range typeMap {
		val, ok := vars[key]
		if !ok {
			continue
		}
		if reason := validate(val, expected); reason != "" {
			violations = append(violations, Violation{
				Key:      key,
				Value:    val,
				Expected: expected,
				Reason:   reason,
			})
		}
	}
	return Result{Violations: violations}, nil
}

func validate(val string, t Type) string {
	switch t {
	case TypeInt:
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return fmt.Sprintf("%q is not a valid integer", val)
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return fmt.Sprintf("%q is not a valid float", val)
		}
	case TypeBool:
		lower := strings.ToLower(val)
		valid := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
		if !valid[lower] {
			return fmt.Sprintf("%q is not a valid bool", val)
		}
	case TypeURL:
		u, err := url.Parse(val)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Sprintf("%q is not a valid URL", val)
		}
	}
	return ""
}

// Format renders a Result as a human-readable string.
func Format(r Result) string {
	if len(r.Violations) == 0 {
		return "all values match declared types"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d type violation(s):\n", len(r.Violations))
	for _, v := range r.Violations {
		fmt.Fprintf(&sb, "  %-24s expected %-6s  got %q  (%s)\n", v.Key, v.Expected, v.Value, v.Reason)
	}
	return strings.TrimRight(sb.String(), "\n")
}
