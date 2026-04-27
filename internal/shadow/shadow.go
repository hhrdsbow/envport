// Package shadow provides functionality to detect keys that exist in one
// snapshot but are shadowed (overridden) by a higher-priority snapshot.
package shadow

import "fmt"

// Result holds information about a shadowed key.
type Result struct {
	Key      string
	BaseVal  string
	OverVal  string
	Source   string // name of the overriding snapshot
}

// Manager is the interface required to load snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Run compares a base snapshot against one or more override snapshots and
// returns keys from base that are shadowed (present and different) in any
// override.
func Run(mgr Manager, baseName string, overrideNames []string) ([]Result, error) {
	baseVars, err := mgr.Load(baseName)
	if err != nil {
		return nil, fmt.Errorf("load base %q: %w", baseName, err)
	}

	var results []Result
	for _, oname := range overrideNames {
		overVars, err := mgr.Load(oname)
		if err != nil {
			return nil, fmt.Errorf("load override %q: %w", oname, err)
		}
		for k, bv := range baseVars {
			ov, ok := overVars[k]
			if ok && ov != bv {
				results = append(results, Result{
					Key:     k,
					BaseVal: bv,
					OverVal: ov,
					Source:  oname,
				})
			}
		}
	}
	return results, nil
}

// Format renders shadow results as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no shadowed keys found\n"
	}
	out := ""
	for _, r := range results {
		out += fmt.Sprintf("%-30s base=%q overridden_by=%s value=%q\n",
			r.Key, r.BaseVal, r.Source, r.OverVal)
	}
	return out
}
