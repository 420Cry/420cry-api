// usercore package provide interface of UserRepository
package usercore

import UserDomain "cry-api/app/domain/users"

type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUUID(uuid string) (*UserDomain.User, error)
	FindByEmail(email string) (*UserDomain.User, error)
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
	FindByVerificationToken(token string) (*UserDomain.User, error)
	FindByAccountVerificationToken(token string) (*UserDomain.User, error)
	FindByUsername(username string) (*UserDomain.User, error)
	FindByUserToken(token string) (*UserDomain.User, error)
	Delete(userID int) error
}
