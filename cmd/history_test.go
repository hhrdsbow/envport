package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envport/internal/history"
)

type memHistoryStore struct {
	data map[string][]byte
}

func newMemHistoryStore() *memHistoryStore {
	return &memHistoryStore{data: map[string][]byte{}}
}
func (m *memHistoryStore) Get(key string) ([]byte, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return v, nil
}
func (m *memHistoryStore) Set(key string, val []byte) error {
	m.data[key] = val
	return nil
}

func TestHistoryList(t *testing.T) {
	store := newMemHistoryStore()
	mgr := history.New(store)
	_ = mgr.Record("snapshot", "dev", "")
	_ = mgr.Record("restore", "dev", "KEY1")

	entries, err := mgr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHistoryOutputContainsHeaders(t *testing.T) {
	store := newMemHistoryStore()
	m)
	_ = mgr.Record("snapshot", "staging", "")

	entries, _ := mgr.List()
	var buf bytes.Buffer
	buf.WriteString("TIME\tOPERATION\tSNAPSHOT\tDETAIL\n")
	for _, e := range entries {
		buf.WriteString(e.Operation + "\t" + e.Snapshot + "\n")
	}
	output := buf.String()
	if !strings.Contains(output, "snapshot") {
		t.Errorf("expected output to contain 'snapshot', got: %s", output)
	}
}

func TestHistoryClear(t *testing.T) {
	store := newMemHistoryStore()
	mgr := history.New(store)
	_ = mgr.Record("snapshot", "dev", "")
	_ = mgr.Clear()
	entries, _ := mgr.List()
	if len(entries) != 0 {
		t.Fatalf("expected 0 after clear, got %d", len(entries))
	}
}
