package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	apiKeyPrefix    = "aif_"
	apiKeyRandSize  = 32
	APIKeyPrefixLen = 12
	argonTime      = 1
	argonMemory    = 64 * 1024
	argonThreads   = 4
	argonKeyLen    = 32
)

// APIKeyPrefix returns the lookup prefix stored in api_keys.key_prefix.
func APIKeyPrefix(plaintext string) (string, error) {
	if !strings.HasPrefix(plaintext, apiKeyPrefix) {
		return "", fmt.Errorf("invalid api key prefix")
	}
	if len(plaintext) < APIKeyPrefixLen {
		return "", fmt.Errorf("api key too short")
	}
	return plaintext[:APIKeyPrefixLen], nil
}

// GenerateAPIKey returns a new plaintext API key and its Argon2id hash.
func GenerateAPIKey() (plaintext, hash string, err error) {
	raw := make([]byte, apiKeyRandSize)
	if _, err = rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("rand: %w", err)
	}

	plaintext = apiKeyPrefix + base64.RawURLEncoding.EncodeToString(raw)
	hash, err = HashAPIKey(plaintext)
	if err != nil {
		return "", "", err
	}
	return plaintext, hash, nil
}

// HashAPIKey hashes an API key with Argon2id (PHC string format).
func HashAPIKey(plaintext string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt: %w", err)
	}

	key := argon2.IDKey([]byte(plaintext), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	)
	return encoded, nil
}

// VerifyAPIKey compares a plaintext key against a stored Argon2id hash.
// Foundation utility only — no persistence or auth flows.
func VerifyAPIKey(plaintext, encoded string) (bool, error) {
	memory, time, threads, salt, expected, err := parseArgon2Hash(encoded)
	if err != nil {
		return false, err
	}

	actual := argon2.IDKey([]byte(plaintext), salt, time, memory, threads, uint32(len(expected)))
	if subtle.ConstantTimeCompare(actual, expected) == 1 {
		return true, nil
	}
	return false, nil
}

func parseArgon2Hash(encoded string) (memory, time uint32, threads uint8, salt, expected []byte, err error) {
	// PHC format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return 0, 0, 0, nil, nil, fmt.Errorf("parse hash: invalid format")
	}

	var version int
	if _, err = fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("parse hash version: %w", err)
	}
	if version != 19 {
		return 0, 0, 0, nil, nil, fmt.Errorf("parse hash: unsupported argon2 version %d", version)
	}

	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("parse hash params: %w", err)
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("decode salt: %w", err)
	}
	expected, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("decode key: %w", err)
	}

	return memory, time, threads, salt, expected, nil
}