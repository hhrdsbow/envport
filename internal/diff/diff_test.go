package diff

import (
	"strings"
	"testing"
)

func TestComputeAdded(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "1", "B": "2"}
	changes := Compute(before, after)
	if len(changes) != 1 || changes[0].Action != "added" || changes[0].Key != "B" {
		t.Fatalf("expected one added change for B, got %+v", changes)
	}
}

func TestComputeRemoved(t *testing.T) {
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "1"}
	changes := Compute(before, after)
	if len(changes) != 1 || changes[0].Action != "removed" || changes[0].Key != "B" {
		t.Fatalf("expected one removed change for B, got %+v", changes)
	}
}

func TestComputeModified(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "2"}
	changes := Compute(before, after)
	if len(changes) != 1 || changes[0].Action != "modified" {
		t.Fatalf("expected one modified change, got %+v", changes)
	}
	if changes[0].Old != "1" || changes[0].New != "2" {
		t.Fatalf("unexpected old/new values: %+v", changes[0])
	}
}

func TestComputeNoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	changes := Compute(env, env)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		change Change
		prefix string
	}{
		{Change{Key: "X", New: "1", Action: "added"}, "+"},
		{Change{Key: "X", Old: "1", Action: "removed"}, "-"},
		{Change{Key: "X", Old: "1", New: "2", Action: "modified"}, "~"},
	}
	for _, tc := range cases {
		out := Format(tc.change)
		if !strings.HasPrefix(out, tc.prefix) {
			t.Errorf("Format(%+v) = %q, want prefix %q", tc.change, out, tc.prefix)
		}
	}
}
