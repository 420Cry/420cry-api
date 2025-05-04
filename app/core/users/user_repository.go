package core

import (
	UserDomain "cry-api/app/domain/users"

	"gorm.io/gorm"
)

// GormUserRepository implements the UserRepository interface for GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository instance
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save persists the user in the database
func (repo *GormUserRepository) Save(user *UserDomain.User) error {
	return repo.db.Create(user).Error
}

// FindByUsernameOrEmail retrieves a user by their username or email
func (repo *GormUserRepository) FindByUsernameOrEmail(username, email string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("username = ?", username).Or("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
