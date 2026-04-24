// Package resolve provides alias-aware snapshot name resolution.
// It transparently maps alias names to their canonical snapshot names.
package resolve

import "fmt"

// AliasStore looks up the canonical name for an alias.
type AliasStore interface {
	Get(alias string) (canonical string, ok bool)
}

// ProfileManager checks whether a snapshot exists by name.
type ProfileManager interface {
	Exists(name string) (bool, error)
}

// Resolver resolves a user-supplied name to a canonical snapshot name.
type Resolver struct {
	aliases  AliasStore
	profiles ProfileManager
}

// New returns a new Resolver.
func New(aliases AliasStore, profiles ProfileManager) *Resolver {
	return &Resolver{aliases: aliases, profiles: profiles}
}

// Resolve returns the canonical snapshot name for the given input.
// If input is an alias, the alias target is returned.
// If input is a direct snapshot name, it is returned as-is.
// An error is returned if neither an alias nor a snapshot is found.
func (r *Resolver) Resolve(input string) (string, error) {
	if canonical, ok := r.aliases.Get(input); ok {
		exists, err := r.profiles.Exists(canonical)
		if err != nil {
			return "", fmt.Errorf("resolve: checking alias target %q: %w", canonical, err)
		}
		if !exists {
			return "", fmt.Errorf("resolve: alias %q points to missing snapshot %q", input, canonical)
		}
		return canonical, nil
	}

	exists, err := r.profiles.Exists(input)
	if err != nil {
		return "", fmt.Errorf("resolve: checking snapshot %q: %w", input, err)
	}
	if !exists {
		return "", fmt.Errorf("resolve: snapshot %q not found", input)
	}
	return input, nil
}

// MustResolve is like Resolve but panics on error. Useful in tests.
func (r *Resolver) MustResolve(input string) string {
	name, err := r.Resolve(input)
	if err != nil {
		panic(err)
	}
	return name
}
