package search

import "strings"

// Snapshot is a minimal interface for searchable snapshots.
type Snapshot interface {
	Name() string
	Tags() []string
	Keys() []string
}

// Query holds search criteria.
type Query struct {
	NameContains string
	Tag          string
	KeyContains  string
}

// Match reports whether s satisfies q.
func Match(s Snapshot, q Query) bool {
	if q.NameContains != "" && !strings.Contains(s.Name(), q.NameContains) {
		return false
	}
	if q.Tag != "" {
		found := false
		for _, t := range s.Tags() {
			if t == q.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if q.KeyContains != "" {
		found := false
		for _, k := range s.Keys() {
			if strings.Contains(k, q.KeyContains) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Filter returns the subset of snapshots matching q.
func Filter(snapshots []Snapshot, q Query) []Snapshot {
	var out []Snapshot
	for _, s := range snapshots {
		if Match(s, q) {
			out = append(out, s)
		}
	}
	return out
}
