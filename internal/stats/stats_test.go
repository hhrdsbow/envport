package stats_test

import (
	"strings"
	"testing"

	"github.com/user/envport/internal/snapshot"
	"github.com/user/envport/internal/stats"
)

type memManager struct {
	data map[string]*snapshot.Snapshot
}

func (m *memManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.data[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return s, nil
}

func seed(vars map[string]string) *snapshot.Snapshot {
	return snapshot.New(vars)
}

func TestRunNoPattern(t *testing.T) {
	mgr := &memManager{
		data: map[string]*snapshot.Snapshot{
			"dev": seed(map[string]string{
				"HOST": "localhost",
				"PORT": "8080",
				"DEBUG": "true",
			}),
		},
	}
	res, err := stats.Run(mgr, "dev", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 {
		t.Errorf("expected 3 keys, got %d", res.Total)
	}
}

func TestRunWithPattern(t *testing.T) {
	mgr := &memManager{
		data: map[string]*snapshot.Snapshot{
			"dev": seed(map[string]string{
				"HOST":     "localhost",
				"DB_HOST":  "db.local",
				"DB_PORT":  "5432",
				"APP_PORT": "8080",
			}),
		},
	}
	res, err := stats.Run(mgr, "dev", "DB_", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 2 {
		t.Errorf("expected 2 matched, got %d", res.Matched)
	}
}

func TestRunTopN(t *testing.T) {
	mgr := &memManager{
		data: map[string]*snapshot.Snapshot{
			"dev": seed(map[string]string{
				"A": "short",
				"B": "a much longer value here",
				"C": "medium length val",
			}),
		},
	}
	res, err := stats.Run(mgr, "dev", "", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.TopByLength) != 2 {
		t.Errorf("expected 2 top entries, got %d", len(res.TopByLength))
	}
	if res.TopByLength[0].Key != "B" {
		t.Errorf("expected B as longest, got %s", res.TopByLength[0].Key)
	}
}

func TestFormat(t *testing.T) {
	res := &stats.Result{
		Name:    "dev",
		Total:   3,
		Matched: 3,
		Empty:   1,
		TopByLength: []stats.KeyLen{
			{Key: "FOO", Length: 10},
		},
	}
	out := stats.Format(res)
	if !strings.Contains(out, "dev") {
		t.Error("expected snapshot name in output")
	}
	if !strings.Contains(out, "Total: 3") {
		t.Error("expected total in output")
	}
	if !strings.Contains(out, "Empty: 1") {
		t.Error("expected empty count in output")
	}
}
