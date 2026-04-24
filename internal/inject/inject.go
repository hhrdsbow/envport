// Package inject applies a snapshot's variables into a running process environment.
package inject

import (
	"fmt"
	"os/exec"

	"github.com/user/envport/internal/snapshot"
)

// Manager provides access to snapshots for injection.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
}

// Options controls injection behaviour.
type Options struct {
	// Keys limits injection to a subset of keys; empty means all keys.
	Keys []string
	// Override replaces existing env vars in the child process when true.
	Override bool
}

// Run executes the given command with the snapshot's variables injected into
// its environment. The parent process environment is inherited first; snapshot
// values are appended (or override, depending on Options.Override).
func Run(mgr Manager, name string, args []string, opts Options) error {
	if len(args) == 0 {
		return fmt.Errorf("inject: no command specified")
	}

	snap, err := mgr.Load(name)
	if err != nil {
		return fmt.Errorf("inject: load snapshot %q: %w", name, err)
	}

	vars := snap.Vars
	if len(opts.Keys) > 0 {
		filtered := make(map[string]string, len(opts.Keys))
		for _, k := range opts.Keys {
			if v, ok := vars[k]; ok {
				filtered[k] = v
			}
		}
		vars = filtered
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = nil

	// Build env slice: inherit current env, then apply snapshot vars.
	cmd.Env = buildEnv(vars, opts.Override)

	cmd.Stdout = nil // callers redirect as needed; keep nil for library use
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("inject: command failed: %w", err)
	}
	return nil
}

// BuildEnv constructs an environment slice suitable for exec.Cmd.Env.
// Exported so cmd layer can use it directly when it needs stdout/stderr wired.
func BuildEnv(vars map[string]string, override bool) []string {
	return buildEnv(vars, override)
}

func buildEnv(vars map[string]string, override bool) []string {
	base := currentEnv()
	if override {
		// Remove existing keys that the snapshot will set.
		filtered := make([]string, 0, len(base))
		for _, e := range base {
			key := envKey(e)
			if _, shadowed := vars[key]; !shadowed {
				filtered = append(filtered, e)
			}
		}
		base = filtered
	}
	for k, v := range vars {
		base = append(base, k+"="+v)
	}
	return base
}

func envKey(pair string) string {
	for i := 0; i < len(pair); i++ {
		if pair[i] == '=' {
			return pair[:i]
		}
	}
	return pair
}
