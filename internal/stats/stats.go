// Package stats provides summary statistics for a named snapshot.
package stats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envport/internal/snapshot"
)

// Manager is the minimal interface required to load a snapshot.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
}

// KeyLen pairs a key name with its value length.
type KeyLen struct {
	Key    string
	Length int
}

// Result holds computed statistics for a snapshot.
type Result struct {
	Name        string
	Total       int
	Matched     int
	Empty       int
	AvgLen      float64
	TopByLength []KeyLen
}

// Run loads the snapshot identified by name and computes statistics.
// pattern filters keys by prefix (empty string matches all).
// topN controls how many longest-value entries to include (0 = all).
func Run(m Manager, name, pattern string, topN int) (*Result, error) {
	snap, err := m.Load(name)
	if err != nil {
		return nil, fmt.Errorf("stats: %w", err)
	}

	vars := snap.Vars
	res := &Result{Name: name, Total: len(vars)}

	var lengths []KeyLen
	totalLen := 0

	for k, v := range vars {
		if pattern != "" && !strings.HasPrefix(k, pattern) {
			continue
		}
		res.Matched++
		if v == "" {
			res.Empty++
		}
		totalLen += len(v)
		lengths = append(lengths, KeyLen{Key: k, Length: len(v)})
	}

	if res.Matched > 0 {
		res.AvgLen = float64(totalLen) / float64(res.Matched)
	}

	sort.Slice(lengths, func(i, j int) bool {
		if lengths[i].Length != lengths[j].Length {
			return lengths[i].Length > lengths[j].Length
		}
		return lengths[i].Key < lengths[j].Key
	})

	if topN > 0 && topN < len(lengths) {
		lengths = lengths[:topN]
	}
	res.TopByLength = lengths

	return res, nil
}

// Format renders a Result as a human-readable summary string.
func Format(r *Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot : %s\n", r.Name)
	fmt.Fprintf(&sb, "Total    : %d\n", r.Total)
	fmt.Fprintf(&sb, "Matched  : %d\n", r.Matched)
	fmt.Fprintf(&sb, "Empty    : %d\n", r.Empty)
	fmt.Fprintf(&sb, "Avg Len  : %.1f\n", r.AvgLen)
	if len(r.TopByLength) > 0 {
		sb.WriteString("Top Keys :\n")
		for _, kl := range r.TopByLength {
			fmt.Fprintf(&sb, "  %-30s %d\n", kl.Key, kl.Length)
		}
	}
	return sb.String()
}
