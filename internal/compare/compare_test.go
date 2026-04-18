package compare_test

import (
	"strings"
	"testing"

	"github.com/envport/envport/internal/compare"
)

type memSnap struct {
	name string
	vars map[string]string
}

func (s *memSnap) Name() string            { return s.name }
func (s *memSnap) Vars() map[string]string { return s.vars }

func TestRunNoChanges(t *testing.T) {
	a := &memSnap{name: "a", vars: map[string]string{"FOO": "bar"}}
	b := &memSnap{name: "b", vars: map[string]string{"FOO": "bar"}}
	r := compare.Run(a, b)
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Modified) != 0 {
		t.Fatal("expected no changes")
	}
	if r.Unchanged["FOO"] != "bar" {
		t.Fatal("expected unchanged FOO")
	}
}

func TestRunAdded(t *testing.T) {
	a := &memSnap{name: "a", vars: map[string]string{}}
	b := &memSnap{name: "b", vars: map[string]string{"NEW": "val"}}
	r := compare.Run(a, b)
	if r.Added["NEW"] != "val" {
		t.Fatal("expected NEW in added")
	}
}

func TestRunRemoved(t *testing.T) {
	a := &memSnap{name: "a", vars: map[string]string{"OLD": "x"}}
	b := &memSnap{name: "b", vars: map[string]string{}}
	r := compare.Run(a, b)
	if r.Removed["OLD"] != "x" {
		t.Fatal("expected OLD in removed")
	}
}

func TestRunModified(t *testing.T) {
	a := &memSnap{name: "a", vars: map[string]string{"K": "v1"}}
	b := &memSnap{name: "b", vars: map[string]string{"K": "v2"}}
	r := compare.Run(a, b)
	if r.Modified["K"] != ([2]string{"v1", "v2"}) {
		t.Fatal("expected K modified")
	}
}

func TestFormat(t *testing.T) {
	a := &memSnap{name: "base", vars: map[string]string{"A": "1", "B": "old"}}
	b := &memSnap{name: "other", vars: map[string]string{"B": "new", "C": "3"}}
	r := compare.Run(a, b)
	out := compare.Format(r)
	if !strings.Contains(out, "compare base..other") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "+ C=3") {
		t.Error("missing added")
	}
	if !strings.Contains(out, "- A=1") {
		t.Error("missing removed")
	}
	if !strings.Contains(out, "~ B: old -> new") {
		t.Error("missing modified")
	}
}
