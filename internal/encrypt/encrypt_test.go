package encrypt_test

import (
	"testing"

	"envport/internal/encrypt"
)

const pass = "s3cr3t"

func TestRoundTrip(t *testing.T) {
	plain := "hello world"
	enc, err := encrypt.Encrypt(plain, pass)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if enc == plain {
		t.Fatal("expected ciphertext to differ from plaintext")
	}
	dec, err := encrypt.Decrypt(enc, pass)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if dec != plain {
		t.Fatalf("got %q, want %q", dec, plain)
	}
}

func TestWrongPassphrase(t *testing.T) {
	enc, _ := encrypt.Encrypt("secret", pass)
	_, err := encrypt.Decrypt(enc, "wrongpass")
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestInvalidBase64(t *testing.T) {
	_, err := encrypt.Decrypt("!!!notbase64!!!", pass)
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestEncryptMapDecryptMap(t *testing.T) {
	vars := map[string]string{"A": "alpha", "B": "beta"}
	enc, err := encrypt.EncryptMap(vars, pass)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	for k, v := range enc {
		if v == vars[k] {
			t.Errorf("key %s not encrypted", k)
		}
	}
	dec, err := encrypt.DecryptMap(enc, pass)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	for k, want := range vars {
		if dec[k] != want {
			t.Errorf("key %s: got %q, want %q", k, dec[k], want)
		}
	}
}

func TestNonDeterministic(t *testing.T) {
	a, _ := encrypt.Encrypt("same", pass)
	b, _ := encrypt.Encrypt("same", pass)
	if a == b {
		t.Fatal("expected different ciphertexts due to random nonce")
	}
}
