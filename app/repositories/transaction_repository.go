package repositorie

import (
	Models "cry-api/app/models"

	"gorm.io/gorm"
)

type GormTransactionRepository struct {
	db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) *GormTransactionRepository {
	return &GormTransactionRepository{db: db}
}

type TransactionRepository interface {
	CreateUserByOAuth(user *Models.User, oauthAccount *Models.Oauth_Accounts) error
}

func (repo *GormTransactionRepository) CreateUserByOAuth(user *Models.User, oauthAccount *Models.Oauth_Accounts) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		var createdUser *Models.User
		err := tx.Where("uuid = ?", user.UUID).First(&createdUser).Error

		if err != nil {
			return err
		}

		oauthAccount.UserId = createdUser.ID

		if err := tx.Create(oauthAccount).Error; err != nil {
			return err
		}

		return nil
	})
}
