package resolve_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envport/internal/resolve"
)

// --- in-memory stubs ---

type memAliases map[string]string

func (m memAliases) Get(alias string) (string, bool) {
	v, ok := m[alias]
	return v, ok
}

type memProfiles map[string]bool

func (m memProfiles) Exists(name string) (bool, error) {
	return m[name], nil
}

type errProfiles struct{}

func (e errProfiles) Exists(_ string) (bool, error) {
	return false, errors.New("storage unavailable")
}

// --- tests ---

func TestResolveDirectSnapshot(t *testing.T) {
	r := resolve.New(memAliases{}, memProfiles{"prod": true})
	got, err := r.Resolve("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "prod" {
		t.Fatalf("want prod, got %s", got)
	}
}

func TestResolveViaAlias(t *testing.T) {
	r := resolve.New(
		memAliases{"live": "prod-2024"},
		memProfiles{"prod-2024": true},
	)
	got, err := r.Resolve("live")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "prod-2024" {
		t.Fatalf("want prod-2024, got %s", got)
	}
}

func TestResolveMissingSnapshot(t *testing.T) {
	r := resolve.New(memAliases{}, memProfiles{})
	_, err := r.Resolve("ghost")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestResolveBrokenAlias(t *testing.T) {
	r := resolve.New(
		memAliases{"broken": "nonexistent"},
		memProfiles{},
	)
	_, err := r.Resolve("broken")
	if err == nil {
		t.Fatal("expected error for alias pointing to missing snapshot")
	}
}

func TestResolveProfileStoreError(t *testing.T) {
	r := resolve.New(memAliases{}, errProfiles{})
	_, err := r.Resolve("anything")
	if err == nil {
		t.Fatal("expected error when profile store fails")
	}
}

func TestResolveAliasStoreError(t *testing.T) {
	r := resolve.New(
		memAliases{"bad": "target"},
		errProfiles{},
	)
	_, err := r.Resolve("bad")
	if err == nil {
		t.Fatal("expected error when profile store fails for alias target")
	}
}
