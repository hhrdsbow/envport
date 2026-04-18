package search_test

import (
	"testing"

	"envport/internal/search"
)

type stub struct {
	name string
	tags []string
	keys []string
}

func (s stub) Name() string   { return s.name }
func (s stub) Tags() []string { return s.tags }
func (s stub) Keys() []string { return s.keys }

var fixtures = []search.Snapshot{
	stub{name: "prod-api", tags: []string{"prod", "api"}, keys: []string{"DATABASE_URL", "PORT"}},
	stub{name: "staging-api", tags: []string{"staging", "api"}, keys: []string{"DATABASE_URL", "DEBUG"}},
	stub{name: "local-worker", tags: []string{"local"}, keys: []string{"REDIS_URL", "PORT"}},
}

func TestMatchName(t *testing.T) {
	q := search.Query{NameContains: "api"}
	results := search.Filter(fixtures, q)
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestMatchTag(t *testing.T) {
	q := search.Query{Tag: "prod"}
	results := search.Filter(fixtures, q)
	if len(results) != 1 || results[0].Name() != "prod-api" {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestMatchKey(t *testing.T) {
	q := search.Query{KeyContains: "REDIS"}
	results := search.Filter(fixtures, q)
	if len(results) != 1 || results[0].Name() != "local-worker" {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestMatchCombined(t *testing.T) {
	q := search.Query{Tag: "api", KeyContains: "DEBUG"}
	results := search.Filter(fixtures, q)
	if len(results) != 1 || results[0].Name() != "staging-api" {
		t.Fatalf("unexpected results: %v", results)
	}
}

func TestMatchNone(t *testing.T) {
	q := search.Query{NameContains: "nonexistent"}
	results := search.Filter(fixtures, q)
	if len(results) != 0 {
		t.Fatalf("expected 0, got %d", len(results))
	}
}

func TestMatchEmpty(t *testing.T) {
	results := search.Filter(fixtures, search.Query{})
	if len(results) != len(fixtures) {
		t.Fatalf("expected all %d, got %d", len(fixtures), len(results))
	}
}
