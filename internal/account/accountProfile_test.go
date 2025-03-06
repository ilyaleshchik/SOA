package acc

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetAccountProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&AccountProfile{})

	expectedProfile := &AccountProfile{AccountID: 1, Name: "aaaaa", LastName: "vvvvv", Bio: "ahahahah", LastUpdateTS: time.Now().Unix()}
	db.Create(expectedProfile)

	profile, err := GetAccountProfile(db, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile == nil || profile.AccountID != expectedProfile.AccountID {
		t.Fatalf("expected profile with ID %d, got %+v", expectedProfile.AccountID, profile)
	}

	profile, err = GetAccountProfile(db, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile != nil {
		t.Fatal("expected nil profile for non-existent account")
	}
}

func TestGetOrCreateAccountProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&AccountProfile{})

	profile, err := GetOrCreateAccountProfile(db, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile == nil || profile.AccountID != 1 {
		t.Fatal("expected new profile with ID 1")
	}

	profile2, err := GetOrCreateAccountProfile(db, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile2 == nil || profile2.AccountID != 1 {
		t.Fatal("expected to retrieve the same profile")
	}
}
