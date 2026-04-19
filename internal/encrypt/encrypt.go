// Package encrypt provides simple AES-GCM encryption for snapshot values.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidCiphertext is returned when decryption fails.
var ErrInvalidCiphertext = errors.New("encrypt: invalid ciphertext")

// deriveKey produces a 32-byte AES key from a passphrase.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// Encrypt encrypts plaintext with the given passphrase using AES-256-GCM.
// The result is base64-encoded (nonce + ciphertext).
func Encrypt(plaintext, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decrypts a base64-encoded ciphertext produced by Encrypt.
func Decrypt(encoded, passphrase string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := aead.NonceSize()
	if len(data) < ns {
		return "", ErrInvalidCiphertext
	}
	plain, err := aead.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plain), nil
}

// EncryptMap encrypts all values in a map, returning a new map.
func EncryptMap(vars map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		enc, err := Encrypt(v, passphrase)
		if err != nil {
			return nil, err
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all values in a map, returning a new map.
func DecryptMap(vars map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		dec, err := Decrypt(v, passphrase)
		if err != nil {
			return nil, err
		}
		out[k] = dec
	}
	return out, nil
}
