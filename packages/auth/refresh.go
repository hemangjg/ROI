package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const refreshTokenBytes = 32

// GenerateRefreshToken returns an opaque refresh token and its SHA-256 hash for storage.
func GenerateRefreshToken() (plaintext, hash string, err error) {
	raw := make([]byte, refreshTokenBytes)
	if _, err = rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("rand: %w", err)
	}

	plaintext = base64.RawURLEncoding.EncodeToString(raw)
	hash = HashRefreshToken(plaintext)
	return plaintext, hash, nil
}

// HashRefreshToken hashes a refresh token for database lookup.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}