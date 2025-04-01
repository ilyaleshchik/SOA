package main

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func generateTestKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	pubKey := &privKey.PublicKey
	return privKey, pubKey, nil
}

func TestGenJWT(t *testing.T) {
	privKey, _, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	s := &Server{jwtPrivate: privKey}
	token, err := s.genJWT("testuser")
	if err != nil {
		t.Fatalf("unexpected error generating JWT: %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty JWT token")
	}
}
