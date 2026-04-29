package count_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/count"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	data map[string]map[string]string
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func seed() *memManager {
	return &memManager{
		data: map[string]map[string]string{
			"prod": {
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5432",
				"APP_SECRET":  "s3cr3t",
				"APP_VERSION": "1.2.3",
				"LOG_LEVEL":   "info",
			},
		},
	}
}

func TestRunNoPattern(t *testing.T) {
	m := seed()
	res, err := count.Run(m, "prod", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 5 {
		t.Errorf("expected Total=5, got %d", res.Total)
	}
	if res.Matched != 5 {
		t.Errorf("expected Matched=5, got %d", res.Matched)
	}
}

func TestRunWithPattern(t *testing.T) {
	m := seed()
	res, err := count.Run(m, "prod", "DB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 5 {
		t.Errorf("expected Total=5, got %d", res.Total)
	}
	if res.Matched != 2 {
		t.Errorf("expected Matched=2, got %d", res.Matched)
	}
}

func TestRunMissing(t *testing.T) {
	m := seed()
	_, err := count.Run(m, "missing", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormatNoPattern(t *testing.T) {
	r := count.Result{Name: "prod", Total: 5, Matched: 5}
	out := count.Format(r)
	if out != "prod: 5 variable(s)" {
		t.Errorf("unexpected format output: %q", out)
	}
}

func TestFormatWithPattern(t *testing.T) {
	r := count.Result{Name: "prod", Total: 5, Matched: 2, Pattern: "DB"}
	out := count.Format(r)
	expected := `prod: 2/5 variable(s) match "DB"`
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
