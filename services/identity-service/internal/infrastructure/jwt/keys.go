package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// Keys holds RSA key material for JWT signing and verification.
type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

// Load reads PEM-encoded RSA private and public keys from disk.
func Load(privatePath, publicPath string) (*Keys, error) {
	privateKey, err := loadRSAPrivateKey(privatePath)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}

	publicKey, err := loadRSAPublicKey(publicPath)
	if err != nil {
		return nil, fmt.Errorf("load public key: %w", err)
	}

	return &Keys{Private: privateKey, Public: publicKey}, nil
}

func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode private key pem")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not rsa")
	}
	return rsaKey, nil
}

func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode public key pem")
	}

	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not rsa")
	}
	return rsaKey, nil
}