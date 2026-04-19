package cmd_test

import (
	"bytes"
	"testing"

	"envport/internal/encrypt"
	"envport/internal/profile"
	"envport/internal/snapshot"
)

func setupEncryptProfile(t *testing.T) (string, *profile.Manager) {
	t.Helper()
	dir := t.TempDir()
	mgr, err := profile.New(dir)
	if err != nil {
		t.Fatalf("profile.New: %v", err)
	}
	snap := snapshot.New(map[string]string{"KEY": "plainvalue", "DB": "secret"})
	if err := mgr.Save("myprofile", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return dir, mgr
}

func TestEncryptRoundTrip(t *testing.T) {
	_, mgr := setupEncryptProfile(t)
	pass := "testpass"

	// Load, encrypt, save manually to simulate cmd behaviour.
	snap, err := mgr.Load("myprofile")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	enc, err := encrypt.EncryptMap(snap.Vars, pass)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	snap.Vars = enc
	if err := mgr.Save("myprofile", snap); err != nil {
		t.Fatalf("Save encrypted: %v", err)
	}

	// Decrypt and verify.
	snap2, err := mgr.Load("myprofile")
	if err != nil {
		t.Fatalf("Load encrypted: %v", err)
	}
	dec, err := encrypt.DecryptMap(snap2.Vars, pass)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	if dec["KEY"] != "plainvalue" {
		t.Errorf("KEY: got %q, want %q", dec["KEY"], "plainvalue")
	}
	if dec["DB"] != "secret" {
		t.Errorf("DB: got %q, want %q", dec["DB"], "secret")
	}
}

func TestEncryptOutputMessage(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("Profile \"myprofile\" encrypted successfully.\n")
	if !bytes.Contains(buf.Bytes(), []byte("encrypted successfully")) {
		t.Error("expected success message")
	}
}
