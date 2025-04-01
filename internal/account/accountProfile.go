package acc

import (
	"time"

	"gorm.io/gorm"
)

type AccountProfile struct {
	AccountID    AccountID `gorm:"uniqueIndex"`
	Bio          string
	Name         string
	LastName     string
	LastUpdateTS int64
}

func GetAccountProfile(db *gorm.DB, accountID AccountID) (*AccountProfile, error) {
	var ap *AccountProfile
	res := db.Model(&AccountProfile{}).Where("account_id = ?", accountID).Find(&ap)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return ap, nil
}

func GetOrCreateAccountProfile(db *gorm.DB, accountID AccountID) (*AccountProfile, error) {
	var ap *AccountProfile
	res := db.Model(&AccountProfile{}).Where("account_id = ?", accountID).Find(&ap)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		ap := &AccountProfile{
			AccountID:    accountID,
			Bio:          "",
			Name:         "",
			LastName:     "",
			LastUpdateTS: time.Now().Unix(),
		}
		err := db.Create(&ap).Error
		if err != nil {
			return nil, err
		}

		return ap, nil
	}

	return ap, nil
}
