package cmd

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"

	"envport/internal/schedule"

	"github.com/spf13/cobra"
)

type memScheduleStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemScheduleStore() *memScheduleStore {
	return &memScheduleStore{data: make(map[string][]byte)}
}
func (s *memScheduleStore) Set(n string, d []byte) error {
	s.mu.Lock(); defer s.mu.Unlock(); s.data[n] = d; return nil
}
func (s *memScheduleStore) Get(n string) ([]byte, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	d, ok := s.data[n]
	if !ok {
		return nil, errors.New("not found")
	}
	return d, nil
}
func (s *memScheduleStore) Delete(n string) error {
	s.mu.Lock(); defer s.mu.Unlock(); delete(s.data, n); return nil
}
func (s *memScheduleStore) List() ([]string, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	var ks []string
	for k := range s.data {
		ks = append(ks, k)
	}
	return ks, nil
}

func buildScheduleCmd(m *schedule.Manager) *cobra.Command {
	root := &cobra.Command{Use: "envport"}

	add := &cobra.Command{
		Use:  "add <name>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, _ := cmd.Flags().GetString("profile")
			interval, _ := cmd.Flags().GetDuration("interval")
			return m.Add(schedule.Entry{Name: args[0], Profile: profile, Interval: interval})
		},
	}
	add.Flags().String("profile", "", "")
	add.Flags().Duration("interval", time.Hour, "")

	list := &cobra.Command{
		Use:  "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := m.List()
			if err != nil {
				return err
			}
			for _, e := range entries {
				cmd.Printf("%s %s %s\n", e.Name, e.Profile, e.Interval)
			}
			return nil
		},
	}

	root.AddCommand(add, list)
	return root
}

func TestScheduleAdd(t *testing.T) {
	m := schedule.New(newMemScheduleStore())
	root := buildScheduleCmd(m)
	root.SetArgs([]string{"add", "nightly", "--profile", "prod", "--interval", "24h"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	e, err := m.Get("nightly")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.Profile != "prod" {
		t.Errorf("profile = %q, want prod", e.Profile)
	}
}

func TestScheduleList(t *testing.T) {
	m := schedule.New(newMemScheduleStore())
	_ = m.Add(schedule.Entry{Name: "daily", Profile: "staging", Interval: 24 * time.Hour})
	root := buildScheduleCmd(m)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("daily")) {
		t.Errorf("output missing 'daily': %s", buf.String())
	}
}
