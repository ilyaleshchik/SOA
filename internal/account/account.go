package acc

import (
	"time"

	"gorm.io/gorm"
)

type AccountID uint32

type Account struct {
	gorm.Model

	ID AccountID `gorm:"primaryKey"`

	Birthday time.Time
	Email    string
	Phone    string
	Username string // login
	Password string
}

func GetAccountByID(db *gorm.DB, accountID AccountID) (*Account, error) {
	var acc *Account
	res := db.Model(&Account{}).Where("id = ?", accountID).Find(&acc)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return acc, nil
}

func GetAccountByUsername(db *gorm.DB, username string) (*Account, error) {
	var acc *Account
	res := db.Model(&Account{}).Where("username = ?", username).Find(&acc)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return acc, nil
}

func SaveAccount(db *gorm.DB, account *Account) error {
	return db.Create(&account).Error
}
