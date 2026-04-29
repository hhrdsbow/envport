// Package stats provides summary statistics for a snapshot's environment variables.
package stats

import (
	"fmt"
	"sort"
	"strings"
)

// Manager describes the minimal interface needed to load a snapshot.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Result holds computed statistics for a snapshot.
type Result struct {
	Name        string
	Total       int
	Empty       int
	Unique      int // values that appear exactly once
	Duplicates  int // values shared by more than one key
	AvgLen      float64
	MaxLen      int
	MinLen      int
	TopPrefixes []PrefixCount
}

// PrefixCount pairs a key prefix with how many keys share it.
type PrefixCount struct {
	Prefix string
	Count  int
}

// Run computes statistics for the named snapshot.
func Run(m Manager, name string) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("stats: load %q: %w", name, err)
	}

	res := Result{Name: name, Total: len(vars)}
	if res.Total == 0 {
		return res, nil
	}

	valueCounts := make(map[string]int, len(vars))
	prefixCounts := make(map[string]int)
	totalLen := 0
	res.MinLen = -1

	for k, v := range vars {
		if v == "" {
			res.Empty++
		}
		valueCounts[v]++
		totalLen += len(v)
		if len(v) > res.MaxLen {
			res.MaxLen = len(v)
		}
		if res.MinLen == -1 || len(v) < res.MinLen {
			res.MinLen = len(v)
		}
		// collect two-segment prefix: e.g. "APP_DB" from "APP_DB_HOST"
		parts := strings.SplitN(k, "_", 3)
		if len(parts) >= 2 {
			prefixCounts[parts[0]+"_"+parts[1]]++
		} else {
			prefixCounts[parts[0]]++
		}
	}

	res.AvgLen = float64(totalLen) / float64(res.Total)

	for _, cnt := range valueCounts {
		if cnt == 1 {
			res.Unique++
		} else {
			res.Duplicates++
		}
	}

	res.TopPrefixes = topN(prefixCounts, 5)
	return res, nil
}

// Format renders a Result as a human-readable summary string.
func Format(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot : %s\n", r.Name)
	fmt.Fprintf(&sb, "Total    : %d\n", r.Total)
	fmt.Fprintf(&sb, "Empty    : %d\n", r.Empty)
	fmt.Fprintf(&sb, "Unique   : %d\n", r.Unique)
	fmt.Fprintf(&sb, "Duplicates: %d\n", r.Duplicates)
	fmt.Fprintf(&sb, "Avg len  : %.1f\n", r.AvgLen)
	fmt.Fprintf(&sb, "Max len  : %d\n", r.MaxLen)
	fmt.Fprintf(&sb, "Min len  : %d\n", r.MinLen)
	if len(r.TopPrefixes) > 0 {
		sb.WriteString("Top prefixes:\n")
		for _, p := range r.TopPrefixes {
			fmt.Fprintf(&sb, "  %-20s %d\n", p.Prefix, p.Count)
		}
	}
	return sb.String()
}

// topN returns up to n prefixes sorted by count descending.
func topN(m map[string]int, n int) []PrefixCount {
	list := make([]PrefixCount, 0, len(m))
	for k, v := range m {
		list = append(list, PrefixCount{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Prefix < list[j].Prefix
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}
