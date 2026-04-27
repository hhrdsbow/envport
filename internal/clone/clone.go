package clone

import (
	"fmt"
	"time"

	"envport/internal/snapshot"
)

// Manager defines the interface required for cloning snapshots.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
	Save(name string, snap *snapshot.Snapshot) error
	List() ([]string, error)
}

// Options configures clone behaviour.
type Options struct {
	Suffix string // appended to dest if dest already exists; defaults to "-copy"
}

// Run clones src snapshot into dest. If dest is empty a name is derived from src.
// The final name used is returned; it may differ from dest if dest was already taken.
func Run(m Manager, src, dest string, opts Options) (string, error) {
	if src == "" {
		return "", fmt.Errorf("source name must not be empty")
	}

	snap, err := m.Load(src)
	if err != nil {
		return "", fmt.Errorf("load %q: %w", src, err)
	}

	if dest == "" {
		dest = src
	}

	suffix := opts.Suffix
	if suffix == "" {
		suffix = "-copy"
	}

	existing, err := m.List()
	if err != nil {
		return "", fmt.Errorf("list snapshots: %w", err)
	}
	taken := make(map[string]bool, len(existing))
	for _, n := range existing {
		taken[n] = true
	}

	final := dest
	for taken[final] {
		final = final + suffix
	}

	cloned := &snapshot.Snapshot{
		Name:      final,
		Vars:      copyVars(snap.Vars),
		CreatedAt: time.Now(),
	}

	if err := m.Save(final, cloned); err != nil {
		return "", fmt.Errorf("save %q: %w", final, err)
	}

	return final, nil
}

func copyVars(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
