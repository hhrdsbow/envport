// Package transform applies key/value transformations to a snapshot's variables.
package transform

import (
	"fmt"
	"strings"
)

// Op represents a single transformation operation.
type Op struct {
	Kind  string // "prefix-add", "prefix-remove", "suffix-add", "suffix-remove", "uppercase", "lowercase"
	Value string // operand (prefix/suffix string), empty for case ops
}

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a transform run.
type Result struct {
	Changed int
	Total   int
}

// Run applies ops to the named snapshot and saves the result.
func Run(m Manager, name string, ops []Op, dryRun bool) (Result, map[string]string, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, nil, fmt.Errorf("transform: load %q: %w", name, err)
	}

	out := make(map[string]string, len(vars))
	changed := 0
	for k, v := range vars {
		newK, newV := k, v
		for _, op := range ops {
			newK, newV = applyOp(op, newK, newV)
		}
		if newK != k || newV != v {
			changed++
		}
		out[newK] = newV
	}

	if !dryRun {
		if err := m.Save(name, out); err != nil {
			return Result{}, nil, fmt.Errorf("transform: save %q: %w", name, err)
		}
	}

	return Result{Changed: changed, Total: len(vars)}, out, nil
}

func applyOp(op Op, k, v string) (string, string) {
	switch op.Kind {
	case "prefix-add":
		return op.Value + k, v
	case "prefix-remove":
		return strings.TrimPrefix(k, op.Value), v
	case "suffix-add":
		return k + op.Value, v
	case "suffix-remove":
		return strings.TrimSuffix(k, op.Value), v
	case "uppercase":
		return strings.ToUpper(k), strings.ToUpper(v)
	case "lowercase":
		return strings.ToLower(k), strings.ToLower(v)
	}
	return k, v
}

// Format returns a human-readable summary of the result.
func Format(r Result, dryRun bool) string {
	action := "transformed"
	if dryRun {
		action = "would transform"
	}
	return fmt.Sprintf("%s %d/%d variable(s)", action, r.Changed, r.Total)
}
