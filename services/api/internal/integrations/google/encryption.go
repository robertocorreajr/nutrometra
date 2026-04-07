package google

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// EncryptToken encrypts plaintext using AES-256-GCM with the given hex-encoded key.
// The hexKey must be 64 hex characters (32 bytes for AES-256).
// Returns the ciphertext as a base64 string (nonce prepended to ciphertext).
func EncryptToken(plaintext string, hexKey string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("encryption: invalid hex key: %w", err)
	}
	if len(key) != 32 {
		return "", fmt.Errorf("encryption: key must be 32 bytes (64 hex chars), got %d bytes", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encryption: new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encryption: new GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encryption: generate nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptToken decrypts a base64-encoded AES-256-GCM ciphertext (with prepended nonce)
// using the given hex-encoded key. Returns the original plaintext.
func DecryptToken(ciphertextB64 string, hexKey string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("decryption: invalid hex key: %w", err)
	}
	if len(key) != 32 {
		return "", fmt.Errorf("decryption: key must be 32 bytes (64 hex chars), got %d bytes", len(key))
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", fmt.Errorf("decryption: invalid base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("decryption: new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("decryption: new GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("decryption: ciphertext too short")
	}

	nonce, encrypted := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("decryption: authentication failed: %w", err)
	}

	return string(plaintext), nil
}
