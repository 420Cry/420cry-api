package repositorie

import (
	Models "cry-api/app/models"

	"gorm.io/gorm"
)

type OAuthRepository interface {
	Save(oauthAccounts *Models.Oauth_Accounts) error
}

type GormOAuthRepository struct {
	db *gorm.DB
}

func NewGormOAuthRepository(db *gorm.DB) *GormOAuthRepository {
	return &GormOAuthRepository{db: db}
}

func (repo *GormOAuthRepository) Save(oauthAccounts *Models.Oauth_Accounts) error {
	return repo.db.Save(oauthAccounts).Error
}
