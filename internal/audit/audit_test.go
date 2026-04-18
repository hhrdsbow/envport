package audit_test

import (
	"testing"
	"time"

	"envport/internal/audit"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	entries []audit.Entry
}

func (m *memStore) Append(e audit.Entry) error {
	m.entries = append(m.entries, e)
	return nil
}

func (m *memStore) List() ([]audit.Entry, error) {
	return m.entries, nil
}

func TestRecordAndList(t *testing.T) {
	st := &memStore{}
	mgr := audit.New(st)

	if err := mgr.Record("snapshot", "prod", "captured 5 vars"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := mgr.Record("restore", "prod", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Operation != "snapshot" {
		t.Errorf("expected snapshot, got %s", entries[0].Operation)
	}
	if entries[1].Detail != "" {
		t.Errorf("expected empty detail, got %s", entries[1].Detail)
	}
}

func TestFormat(t *testing.T) {
	e := audit.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Operation: "copy",
		Profile:   "staging",
		Detail:    "-> staging-backup",
	}
	out := audit.Format(e)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	for _, substr := range []string{"copy", "staging", "staging-backup"} {
		if !contains(out, substr) {
			t.Errorf("expected %q in output %q", substr, out)
		}
	}
}

func TestFormatNoDetail(t *testing.T) {
	e := audit.Entry{
		Timestamp: time.Now().UTC(),
		Operation: "delete",
		Profile:   "old",
	}
	out := audit.Format(e)
	if !contains(out, "delete") || !contains(out, "old") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
