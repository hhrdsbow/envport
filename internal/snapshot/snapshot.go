package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents a saved set of environment variables.
type Snapshot struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot from the current environment.
func New(name string, environ []string) *Snapshot {
	env := make(map[string]string, len(environ))
	for _, e := range environ {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				env[e[:i]] = e[i+1:]
				break
			}
		}
	}
	return &Snapshot{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Env:       env,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// ToExports returns a slice of "export KEY=VALUE" strings.
func (s *Snapshot) ToExports() []string {
	lines := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		lines = append(lines, "export "+k+"="+v)
	}
	return lines
}
