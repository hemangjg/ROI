package auth_test

import (
	"testing"

	"github.com/ai-finops/ai-finops/packages/auth"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := auth.HashPassword("demo-password-change-me")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if !auth.VerifyPassword("demo-password-change-me", hash) {
		t.Fatal("expected password to verify")
	}
	if auth.VerifyPassword("wrong-password", hash) {
		t.Fatal("expected wrong password to fail")
	}
}