// Package prefix provides operations to add or remove key prefixes across snapshots.
package prefix

import (
	"fmt"
	"strings"
)

// Manager describes the snapshot operations required by this package.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a prefix operation.
type Result struct {
	Added   []string
	Removed []string
	Skipped []string
}

// RunAdd adds prefix p to every key in the snapshot named src, saving the
// result to dst. If src == dst the snapshot is updated in place.
// Keys that already start with p are kept unchanged and recorded in
// Result.Skipped to avoid double-prefixing.
func RunAdd(m Manager, src, dst, p string) (Result, error) {
	if p == "" {
		return Result{}, fmt.Errorf("prefix must not be empty")
	}
	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("load %q: %w", src, err)
	}
	out := make(map[string]string, len(vars))
	var res Result
	for k, v := range vars {
		if strings.HasPrefix(k, p) {
			out[k] = v
			res.Skipped = append(res.Skipped, k)
		} else {
			nk := p + k
			out[nk] = v
			res.Added = append(res.Added, nk)
		}
	}
	if err := m.Save(dst, out); err != nil {
		return Result{}, fmt.Errorf("save %q: %w", dst, err)
	}
	return res, nil
}

// RunRemove strips prefix p from every key in src that carries it, saving to
// dst. Keys that do not start with p are kept unchanged and recorded in
// Result.Skipped.
func RunRemove(m Manager, src, dst, p string) (Result, error) {
	if p == "" {
		return Result{}, fmt.Errorf("prefix must not be empty")
	}
	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("load %q: %w", src, err)
	}
	out := make(map[string]string, len(vars))
	var res Result
	for k, v := range vars {
		if strings.HasPrefix(k, p) {
			nk := strings.TrimPrefix(k, p)
			out[nk] = v
			res.Removed = append(res.Removed, k)
		} else {
			out[k] = v
			res.Skipped = append(res.Skipped, k)
		}
	}
	if err := m.Save(dst, out); err != nil {
		return Result{}, fmt.Errorf("save %q: %w", dst, err)
	}
	return res, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("modified : %d\n", len(r.Added)+len(r.Removed)))
	sb.WriteString(fmt.Sprintf("skipped  : %d\n", len(r.Skipped)))
	return sb.String()
}
