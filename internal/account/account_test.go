package acc

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetAccountByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Account{})

	expectedAccount := &Account{ID: 1, Username: "testuser"}
	db.Create(expectedAccount)

	account, err := GetAccountByID(db, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account == nil || account.ID != expectedAccount.ID {
		t.Fatalf("expected account with ID %d, got %+v", expectedAccount.ID, account)
	}

	account, err = GetAccountByID(db, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account != nil {
		t.Fatal("expected nil account for non-existent ID")
	}
}

func TestGetAccountByUsername(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Account{})

	expectedAccount := &Account{ID: 1, Username: "testuser"}
	db.Create(expectedAccount)

	account, err := GetAccountByUsername(db, "testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account == nil || account.Username != expectedAccount.Username {
		t.Fatalf("expected account with username %s, got %+v", expectedAccount.Username, account)
	}

	account, err = GetAccountByUsername(db, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account != nil {
		t.Fatal("expected nil account for non-existent username")
	}
}

func TestSaveAccount(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Account{})

	newAccount := &Account{ID: 1, Username: "newuser"}
	err = SaveAccount(db, newAccount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	savedAccount, err := GetAccountByID(db, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedAccount == nil || savedAccount.ID != newAccount.ID {
		t.Fatal("expected saved account to exist in database")
	}
}
