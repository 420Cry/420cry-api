// cry-api/app/infrastructure/persistence/gorm_user_repository.go
package persistence

import (
	"fmt"
	"time"

	UserDomain "cry-api/app/domain/users"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (repo *GormUserRepository) Save(user *UserDomain.User) error {
	return repo.db.Save(user).Error
}

func (repo *GormUserRepository) FindByUUID(uuid string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *GormUserRepository) FindByEmail(email string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *GormUserRepository) FindByUsernameOrEmail(username string, email string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *GormUserRepository) FindByUsername(username string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *GormUserRepository) FindByVerificationToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("verification_tokens = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if time.Since(user.VerificationTokenCreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("verification token expired")
	}
	return &user, nil
}

func (repo *GormUserRepository) FindByUserToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("verification_tokens = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if time.Since(user.VerificationTokenCreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("user token expired")
	}
	return &user, nil
}

func (repo *GormUserRepository) FindByAccountVerificationToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *GormUserRepository) Delete(userID int) error {
	if err := repo.db.Delete(&UserDomain.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
