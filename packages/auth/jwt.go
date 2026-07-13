package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims are foundation-level JWT claims. No auth flows yet.
type JWTClaims struct {
	UserID string   `json:"uid"`
	OrgID  string   `json:"oid"`
	Roles  []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// SignJWT creates an RS256-signed JWT.
func SignJWT(privateKey any, claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return signed, nil
}

// ParseJWT validates and parses an RS256 JWT.
func ParseJWT(publicKey any, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// NewAccessClaims builds short-lived access token claims.
func NewAccessClaims(userID, orgID string, roles []string, ttl time.Duration) JWTClaims {
	now := time.Now()
	return JWTClaims{
		UserID: userID,
		OrgID:  orgID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
}