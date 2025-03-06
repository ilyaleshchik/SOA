package main

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	acc "soa-hw-ilyaleshchyk/internal/account"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func TestCheckJwt(t *testing.T) {
	privKey, pubKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&acc.Account{})
	db.Create(&acc.Account{ID: 1, Username: "testuser"})

	s := &Server{jwtPrivate: privKey, jwtPublic: pubKey, db: db}
	token, _ := s.genJWT("testuser")

	r := gin.Default()
	r.GET("/test", s.checkJwt, func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Auth", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Auth", "invalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for invalid token, got %d", w.Code)
	}
}
