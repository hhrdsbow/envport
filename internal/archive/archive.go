// Package archive provides functionality to bundle multiple snapshots
// into a single portable archive file and restore them from it.
package archive

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/user/envport/internal/snapshot"
)

// Manager defines the interface for loading and saving snapshots.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
	Save(name string, snap *snapshot.Snapshot) error
	List() ([]string, error)
}

// Archive holds a collection of named snapshots.
type Archive struct {
	CreatedAt time.Time                      `json:"created_at"`
	Snapshots map[string]*snapshot.Snapshot  `json:"snapshots"`
}

// Pack collects the given snapshot names from the manager and returns an Archive.
func Pack(m Manager, names []string) (*Archive, error) {
	ar := &Archive{
		CreatedAt: time.Now().UTC(),
		Snapshots: make(map[string]*snapshot.Snapshot, len(names)),
	}
	for _, name := range names {
		snap, err := m.Load(name)
		if err != nil {
			return nil, fmt.Errorf("archive: loading %q: %w", name, err)
		}
		ar.Snapshots[name] = snap
	}
	return ar, nil
}

// Unpack restores all snapshots contained in the archive via the manager.
// Existing snapshots with the same name are overwritten.
func Unpack(m Manager, ar *Archive) ([]string, error) {
	var restored []string
	for name, snap := range ar.Snapshots {
		if err := m.Save(name, snap); err != nil {
			return restored, fmt.Errorf("archive: saving %q: %w", name, err)
		}
		restored = append(restored, name)
	}
	return restored, nil
}

// Marshal serialises the archive to JSON bytes.
func Marshal(ar *Archive) ([]byte, error) {
	return json.MarshalIndent(ar, "", "  ")
}

// Unmarshal deserialises an archive from JSON bytes.
func Unmarshal(data []byte) (*Archive, error) {
	var ar Archive
	if err := json.Unmarshal(data, &ar); err != nil {
		return nil, fmt.Errorf("archive: unmarshal: %w", err)
	}
	if ar.Snapshots == nil {
		ar.Snapshots = make(map[string]*snapshot.Snapshot)
	}
	return &ar, nil
}
