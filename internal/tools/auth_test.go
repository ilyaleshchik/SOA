package tools

import (
	acc "soa-hw-ilyaleshchyk/internal/account"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Fatal("expected hashed password, got empty string")
	}
}

func TestIsSamePass(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = IsSamePass([]byte(password), hashedPassword)
	if err != nil {
		t.Fatalf("password should match, but got error: %v", err)
	}

	wrongPassword := "wrongpassword"
	err = IsSamePass([]byte(wrongPassword), hashedPassword)
	if err == nil {
		t.Fatal("expected error for incorrect password, but got nil")
	}
}

func TestGenerateRandomAccountID(t *testing.T) {

	generatedIDs := make(map[acc.AccountID]struct{})

	for range 10000 {
		id1 := GenerateRandomAccountID()
		id2 := GenerateRandomAccountID()
		if id1 == id2 {
			t.Logf("Warning: Generated the same ID twice (%d), which is unlikely", id1)
		}

		_, ok := generatedIDs[id1]
		if ok {
			t.Logf("Warning: Generated the same ID twice (%d), which is unlikely", id1)
		}

		_, ok = generatedIDs[id2]
		if ok {
			t.Logf("Warning: Generated the same ID twice (%d), which is unlikely", id2)
		}
		generatedIDs[id1] = struct{}{}
		generatedIDs[id2] = struct{}{}
	}
}
