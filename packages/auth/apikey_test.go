package auth_test

import (
	"strings"
	"testing"

	"github.com/ai-finops/ai-finops/packages/auth"
)

func TestAPIKeyHashRoundTrip(t *testing.T) {
	plaintext, hash, err := auth.GenerateAPIKey()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if !strings.HasPrefix(plaintext, "aif_") {
		t.Fatalf("prefix missing: %q", plaintext)
	}

	prefix, err := auth.APIKeyPrefix(plaintext)
	if err != nil {
		t.Fatalf("prefix: %v", err)
	}
	if len(prefix) != auth.APIKeyPrefixLen {
		t.Fatalf("prefix len = %d, want %d", len(prefix), auth.APIKeyPrefixLen)
	}

	ok, err := auth.VerifyAPIKey(plaintext, hash)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatal("expected valid key")
	}

	ok, err = auth.VerifyAPIKey("aif_invalid", hash)
	if err != nil {
		t.Fatalf("verify invalid: %v", err)
	}
	if ok {
		t.Fatal("expected invalid key")
	}
}

func TestHashAPIKeyDeterministicVerify(t *testing.T) {
	plaintext := "aif_test_key_for_hash_roundtrip"
	hash, err := auth.HashAPIKey(plaintext)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	ok, err := auth.VerifyAPIKey(plaintext, hash)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatal("expected hash to verify")
	}
}