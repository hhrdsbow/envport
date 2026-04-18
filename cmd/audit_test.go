package cmd

import (
	"bytes"
	"testing"
	"time"

	"envport/internal/audit"
)

// memAuditStore is an in-memory audit store for testing.
type memAuditStore struct {
	entries []audit.Entry
}

func (m *memAuditStore) Append(e audit.Entry) error {
	m.entries = append(m.entries, e)
	return nil
}

func (m *memAuditStore) List() ([]audit.Entry, error) {
	return m.entries, nil
}

func TestAuditListOutput(t *testing.T) {
	st := &memAuditStore{}
	mgr := audit.New(st)
	_ = mgr.Record("snapshot", "prod", "3 vars")
	_ = mgr.Record("restore", "dev", "")

	entries, _ := mgr.List()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	var buf bytes.Buffer
	for _, e := range entries {
		buf.WriteString(audit.Format(e))
		buf.WriteString("\n")
	}
	out := buf.String()
	for _, want := range []string{"snapshot", "prod", "restore", "dev"} {
		if !bytesContains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestAuditTimestamp(t *testing.T) {
	e := audit.Entry{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Operation: "merge",
		Profile:   "base",
	}
	out := audit.Format(e)
	if !bytesContains(out, "2024-06-01") {
		t.Errorf("expected date in output: %s", out)
	}
}

func bytesContains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}
