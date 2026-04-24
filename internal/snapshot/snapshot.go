package snapshot

import (
	"encoding/json"
	"fmt"
	"time"
)

// Snapshot holds a named set of environment variables captured at a point in time.
type Snapshot struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// Manager describes the persistence operations a snapshot store must support.
type Manager interface {
	Save(s *Snapshot) error
	Load(name string) (*Snapshot, error)
	Delete(name string) error
	List() ([]string, error)
}

// New creates a new Snapshot with the given name and variable map.
func New(name string, vars map[string]string) *Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Snapshot{
		Name:      name,
		Vars:      copy,
		CreatedAt: time.Now().UTC(),
	}
}

// Load deserialises a Snapshot from JSON bytes.
func Load(data []byte) (*Snapshot, error) {
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Bytes serialises the Snapshot to JSON.
func (s *Snapshot) Bytes() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}
	return b, nil
}
