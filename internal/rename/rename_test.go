package rename_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/rename"
)

type memManager struct {
	data map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *memManager) Save(name string, env map[string]string) error {
	m.data[name] = env
	return nil
}

func (m *memManager) Delete(name string) error {
	delete(m.data, name)
	return nil
}

func TestRunRenamesSnapshot(t *testing.T) {
	mgr := newMemManager()
	_ = mgr.Save("old", map[string]string{"FOO": "bar"})

	if err := rename.Run(mgr, "old", "new"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := mgr.Load("new"); err != nil {
		t.Error("expected 'new' snapshot to exist")
	}
	if _, err := mgr.Load("old"); err == nil {
		t.Error("expected 'old' snapshot to be deleted")
	}
}

func TestRunMissingSource(t *testing.T) {
	mgr := newMemManager()
	err := rename.Run(mgr, "missing", "dst")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
	var nf *rename.ErrNotFound
	if !errors.As(err, &nf) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestRunSameName(t *testing.T) {
	mgr := newMemManager()
	_ = mgr.Save("snap", map[string]string{})
	err := rename.Run(mgr, "snap", "snap")
	if err == nil {
		t.Fatal("expected error when src == dst")
	}
}
