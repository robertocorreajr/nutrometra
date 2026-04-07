package google

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestKey(t *testing.T) string {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return hex.EncodeToString(key)
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	hexKey := generateTestKey(t)
	tests := []struct {
		name      string
		plaintext string
	}{
		{"simple token", "ya29.a0ARrdaM-abc123"},
		{"empty string", ""},
		{"long token", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0"},
		{"special characters", "token/with+special=chars&more%20stuff"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptToken(tt.plaintext, hexKey)
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)
			assert.NotEqual(t, tt.plaintext, encrypted)

			decrypted, err := DecryptToken(encrypted, hexKey)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestEncryptDifferentCiphertexts(t *testing.T) {
	hexKey := generateTestKey(t)
	plaintext := "same-token-value"

	enc1, err := EncryptToken(plaintext, hexKey)
	require.NoError(t, err)
	enc2, err := EncryptToken(plaintext, hexKey)
	require.NoError(t, err)

	// Due to random nonce, same plaintext should produce different ciphertexts.
	assert.NotEqual(t, enc1, enc2)

	// Both should decrypt to the same value.
	dec1, err := DecryptToken(enc1, hexKey)
	require.NoError(t, err)
	dec2, err := DecryptToken(enc2, hexKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, dec1)
	assert.Equal(t, plaintext, dec2)
}

func TestEncryptWithInvalidKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"too short", "abcd1234"},
		{"not hex", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"},
		{"empty key", ""},
		{"odd length hex", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncryptToken("test-token", tt.key)
			assert.Error(t, err)
		})
	}
}

func TestDecryptWithInvalidKey(t *testing.T) {
	goodKey := generateTestKey(t)
	wrongKey := generateTestKey(t)

	encrypted, err := EncryptToken("my-secret-token", goodKey)
	require.NoError(t, err)

	// Wrong key should fail authentication.
	_, err = DecryptToken(encrypted, wrongKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	hexKey := generateTestKey(t)
	encrypted, err := EncryptToken("my-secret-token", hexKey)
	require.NoError(t, err)

	// Decode, tamper, re-encode.
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	require.NoError(t, err)
	require.True(t, len(raw) > 5)

	// Flip a byte in the ciphertext (not the nonce).
	raw[len(raw)-3] ^= 0xFF
	tampered := base64.StdEncoding.EncodeToString(raw)

	_, err = DecryptToken(tampered, hexKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestDecryptInvalidBase64(t *testing.T) {
	hexKey := generateTestKey(t)
	_, err := DecryptToken("not-valid-base64!!!", hexKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid base64")
}

func TestDecryptTooShortCiphertext(t *testing.T) {
	hexKey := generateTestKey(t)
	short := base64.StdEncoding.EncodeToString([]byte("abc"))
	_, err := DecryptToken(short, hexKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}
