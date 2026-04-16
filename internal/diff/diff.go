package diff

import "fmt"

// Change represents a single environment variable change.
type Change struct {
	Key    string
	Old    string
	New    string
	Action string // "added", "removed", "modified"
}

// Compute returns the list of changes between two env maps.
func Compute(before, after map[string]string) []Change {
	var changes []Change

	for k, newVal := range after {
		oldVal, exists := before[k]
		if !exists {
			changes = append(changes, Change{Key: k, Old: "", New: newVal, Action: "added"})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Old: oldVal, New: newVal, Action: "modified"})
		}
	}

	for k, oldVal := range before {
		if _, exists := after[k]; !exists {
			changes = append(changes, Change{Key: k, Old: oldVal, New: "", Action: "removed"})
		}
	}

	return changes
}

// Format returns a human-readable summary of a Change.
func Format(c Change) string {
	switch c.Action {
	case "added":
		return fmt.Sprintf("+ %s=%s", c.Key, c.New)
	case "removed":
		return fmt.Sprintf("- %s=%s", c.Key, c.Old)
	case "modified":
		return fmt.Sprintf("~ %s: %s -> %s", c.Key, c.Old, c.New)
	}
	return ""
}
