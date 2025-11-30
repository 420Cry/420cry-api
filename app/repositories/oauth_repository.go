package repositorie

import (
	Models "cry-api/app/models"

	"gorm.io/gorm"
)

type OAuthRepository interface {
	Save(oauthAccounts *Models.Oauth_Accounts) error
	FindByProviderAndId(provider, providerId string) (*Models.Oauth_Accounts, error)
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

func (repo *GormOAuthRepository) FindByProviderAndId(provider, providerId string) (*Models.Oauth_Accounts, error) {
	var oauthAccount *Models.Oauth_Accounts
	err := repo.db.Where("provider = ?", provider).Where("provider_id = ?", providerId).First(&oauthAccount).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return oauthAccount, nil
}
