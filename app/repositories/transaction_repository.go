package repositorie

import (
	Models "cry-api/app/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

type TransactionRepositoryInterface interface {
	CreateUserByGoogleAuth(user *Models.User, oauthAccount *Models.Oauth_Accounts) error
}

func (repo *TransactionRepository) CreateUserByGoogleAuth(user *Models.User, oauthAccount *Models.Oauth_Accounts) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Create(oauthAccount).Error; err != nil {
			return err
		}

		return nil
	})
}
