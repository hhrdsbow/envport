package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envport/internal/export"
	"envport/internal/profile"
	"envport/internal/snapshot"
)

func setupExportProfile(t *testing.T, dir string) profile.Manager {
	t.Helper()
	mgr, err := profile.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	snap := snapshot.New("testprofile", map[string]string{"KEY": "value", "ANOTHER": "thing"})
	if err := mgr.Save(snap); err != nil {
		t.Fatal(err)
	}
	return mgr
}

func TestExportShell(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var b bytes.Buffer
	if err := export.Write(&b, env, export.FormatShell); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "export KEY") {
		t.Errorf("expected shell export, got: %s", b.String())
	}
}

func TestExportDotenv(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var b bytes.Buffer
	if err := export.Write(&b, env, export.FormatDotenv); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "KEY=value") {
		t.Errorf("expected dotenv line, got: %s", b.String())
	}
}

func TestExportJSON(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var b bytes.Buffer
	if err := export.Write(&b, env, export.FormatJSON); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, "{") || !strings.Contains(out, `"KEY"`) {
		t.Errorf("expected json output, got: %s", out)
	}
}
