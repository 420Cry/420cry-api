package repositorie_test

import (
	"fmt"
	"regexp"
	"testing"

	"cry-api/app/models"
	repositorie "cry-api/app/repositories"
	mocks "cry-api/tests/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormUserRepository_Save_Update(t *testing.T) {
	//revive:disable:redefines-builtin-id
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositorie.NewGormUserRepository(db)

	user := &models.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserRepository_FindByUUID(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositorie.NewGormUserRepository(db)

	user := models.User{
		ID:       1,
		UUID:     "uuid-1234",
		Username: "johndoe",
		Email:    "john@example.com",
	}

	rows := sqlmock.NewRows([]string{"id", "uuid", "username", "email"}).
		AddRow(user.ID, user.UUID, user.Username, user.Email)

	// Success case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(user.UUID, 1).
		WillReturnRows(rows)

	result, err := repo.FindByUUID(user.UUID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.UUID, result.UUID)

	// Not found case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("nonexistent", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err = repo.FindByUUID("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// DB error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("error", 1).
		WillReturnError(fmt.Errorf("db error"))

	result, err = repo.FindByUUID("error")
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserRepository_FindByEmail(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositorie.NewGormUserRepository(db)

	email := "test@example.com"
	user := models.User{
		ID:       1,
		UUID:     "uuid-123",
		Username: "testuser",
		Email:    email,
	}

	rows := sqlmock.NewRows([]string{"id", "uuid", "username", "email"}).
		AddRow(user.ID, user.UUID, user.Username, user.Email)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	result, err := repo.FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)

	// Not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("missing@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err = repo.FindByEmail("missing@example.com")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// DB error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("error@example.com", 1).
		WillReturnError(fmt.Errorf("db error"))

	result, err = repo.FindByEmail("error@example.com")
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserRepository_FindByUsernameOrEmail(t *testing.T) {
	db, mock, close := mocks.SetupMockDB(t)
	defer close()

	repo := repositorie.NewGormUserRepository(db)

	username := "user1"
	email := "user1@example.com"
	user := models.User{
		ID:       1,
		UUID:     "uuid-1",
		Username: username,
		Email:    email,
	}

	rows := sqlmock.NewRows([]string{"id", "uuid", "username", "email"}).
		AddRow(user.ID, user.UUID, user.Username, user.Email)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 OR email = $2 ORDER BY "users"."id" LIMIT $3`)).
		WithArgs(username, email, 1).
		WillReturnRows(rows)

	result, err := repo.FindByUsernameOrEmail(username, email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, username, result.Username)

	// Not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 OR email = $2 ORDER BY "users"."id" LIMIT $3`)).
		WithArgs("missing", "missing@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err = repo.FindByUsernameOrEmail("missing", "missing@example.com")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// DB error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 OR email = $2 ORDER BY "users"."id" LIMIT $3`)).
		WithArgs("error", "error@example.com", 1).
		WillReturnError(fmt.Errorf("db error"))

	result, err = repo.FindByUsernameOrEmail("error", "error@example.com")
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserRepository_FindByUsername(t *testing.T) {
	db, mock, close := mocks.SetupMockDB(t)
	defer close()

	repo := repositorie.NewGormUserRepository(db)

	username := "user1"
	user := models.User{
		ID:       1,
		UUID:     "uuid-1",
		Username: username,
		Email:    "email@example.com",
	}

	rows := sqlmock.NewRows([]string{"id", "uuid", "username", "email"}).
		AddRow(user.ID, user.UUID, user.Username, user.Email)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(username, 1).
		WillReturnRows(rows)

	result, err := repo.FindByUsername(username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, username, result.Username)

	// Not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("missing", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err = repo.FindByUsername("missing")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// DB error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("error", 1).
		WillReturnError(fmt.Errorf("db error"))

	result, err = repo.FindByUsername("error")
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserRepository_Delete(t *testing.T) {
	db, mock, close := mocks.SetupMockDB(t)
	defer close()

	repo := repositorie.NewGormUserRepository(db)

	userID := 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(userID)
	assert.NoError(t, err)

	// Failure case
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("delete error"))
	mock.ExpectRollback()

	err = repo.Delete(userID)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
