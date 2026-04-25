// Package swap provides functionality to exchange variable values between two snapshots.
package swap

import (
	"fmt"
)

// Manager describes what swap needs from a profile store.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a swap operation.
type Result struct {
	SrcName string
	DstName string
	Keys     []string
}

// Run swaps the values of the given keys between src and dst snapshots.
// If keys is empty, all keys present in both snapshots are swapped.
func Run(m Manager, src, dst string, keys []string) (Result, error) {
	srcVars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("swap: load %q: %w", src, err)
	}
	dstVars, err := m.Load(dst)
	if err != nil {
		return Result{}, fmt.Errorf("swap: load %q: %w", dst, err)
	}

	swapKeys := keys
	if len(swapKeys) == 0 {
		swapKeys = commonKeys(srcVars, dstVars)
	}
	if len(swapKeys) == 0 {
		return Result{SrcName: src, DstName: dst}, nil
	}

	newSrc := copyMap(srcVars)
	newDst := copyMap(dstVars)

	for _, k := range swapKeys {
		newSrc[k], newDst[k] = dstVars[k], srcVars[k]
	}

	if err := m.Save(src, newSrc); err != nil {
		return Result{}, fmt.Errorf("swap: save %q: %w", src, err)
	}
	if err := m.Save(dst, newDst); err != nil {
		return Result{}, fmt.Errorf("swap: save %q: %w", dst, err)
	}

	return Result{SrcName: src, DstName: dst, Keys: swapKeys}, nil
}

// Format returns a human-readable summary of the swap result.
func Format(r Result) string {
	if len(r.Keys) == 0 {
		return fmt.Sprintf("no common keys found between %q and %q", r.SrcName, r.DstName)
	}
	out := fmt.Sprintf("swapped %d key(s) between %q and %q:\n", len(r.Keys), r.SrcName, r.DstName)
	for _, k := range r.Keys {
		out += fmt.Sprintf("  %s\n", k)
	}
	return out
}

func commonKeys(a, b map[string]string) []string {
	var keys []string
	for k := range a {
		if _, ok := b[k]; ok {
			keys = append(keys, k)
		}
	}
	return keys
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
