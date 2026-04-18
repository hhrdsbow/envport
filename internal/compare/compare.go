// Package compare provides snapshot comparison utilities.
package compare

import (
	"fmt"
	"sort"
	"strings"
)

// Snapshot is a minimal interface for snapshots needed by compare.
type Snapshot interface {
	Name() string
	Vars() map[string]string
}

// Result holds the comparison between two snapshots.
type Result struct {
	Base  string
	Other string
	Added    map[string]string
	Removed  map[string]string
	Modified map[string][2]string // key -> [baseVal, otherVal]
	Unchanged map[string]string
}

// Run compares two snapshots and returns a Result.
func Run(base, other Snapshot) Result {
	r := Result{
		Base:      base.Name(),
		Other:     other.Name(),
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Modified:  make(map[string][2]string),
		Unchanged: make(map[string]string),
	}
	bv := base.Vars()
	ov := other.Vars()
	for k, v := range bv {
		if nv, ok := ov[k]; !ok {
			r.Removed[k] = v
		} else if nv != v {
			r.Modified[k] = [2]string{v, nv}
		} else {
			r.Unchanged[k] = v
		}
	}
	for k, v := range ov {
		if _, ok := bv[k]; !ok {
			r.Added[k] = v
		}
	}
	return r
}

// Format returns a human-readable summary of the Result.
func Format(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "compare %s..%s\n", r.Base, r.Other)
	for _, k := range sortedKeys(r.Added) {
		fmt.Fprintf(&sb, "+ %s=%s\n", k, r.Added[k])
	}
	for _, k := range sortedKeys(r.Removed) {
		fmt.Fprintf(&sb, "- %s=%s\n", k, r.Removed[k])
	}
	for _, k := range sortedModified(r.Modified) {
		v := r.Modified[k]
		fmt.Fprintf(&sb, "~ %s: %s -> %s\n", k, v[0], v[1])
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedModified(m map[string][2]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
