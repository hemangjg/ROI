package auth_test

import (
	"testing"

	"github.com/ai-finops/ai-finops/packages/auth"
)

func TestRefreshTokenHashDeterministic(t *testing.T) {
	plaintext, hash, err := auth.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if hash != auth.HashRefreshToken(plaintext) {
		t.Fatal("expected deterministic refresh token hash")
	}
}