package tests

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	UserModel "cry-api/app/models"
	repositorie "cry-api/app/repositories"
	mocks "cry-api/tests/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormUserTokenRepository_Save(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := repositorie.NewGormUserTokenRepository(db)

	token := &UserModel.UserToken{
		ID:        0, // force INSERT
		UserID:    1,
		Token:     "abc123",
		Purpose:   "test",
		Consumed:  false,
		ExpiresAt: time.Now().Add(time.Hour),
		UsedAt:    nil,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "user_tokens" ("user_id","token","purpose","expires_at","consumed","used_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "created_at","id"`)).
		WithArgs(token.UserID, token.Token, token.Purpose, sqlmock.AnyArg(), token.Consumed, nil).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "id"}).AddRow(time.Now(), 1))
	mock.ExpectCommit()

	err := repo.Save(token)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserTokenRepository_FindValidToken(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()
	repo := repositorie.NewGormUserTokenRepository(db)

	now := time.Now()
	token := &UserModel.UserToken{
		ID:        1,
		UserID:    1,
		Token:     "abc123",
		Purpose:   "reset_password",
		Consumed:  false,
		ExpiresAt: now.Add(time.Hour),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "token", "purpose", "consumed", "expires_at", "used_at"}).
		AddRow(token.ID, token.UserID, token.Token, token.Purpose, token.Consumed, token.ExpiresAt, nil)

	// GORM passes the LIMIT value as the last argument automatically, so we match any value
	query := `SELECT \* FROM "user_tokens" WHERE token = \$1 AND purpose = \$2 AND consumed = \$3 AND used_at IS NULL AND expires_at > \$4 ORDER BY "user_tokens"\."id"`

	// Success case
	mock.ExpectQuery(query).
		WithArgs(token.Token, token.Purpose, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	result, err := repo.FindValidToken(token.Token, token.Purpose)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token.Token, result.Token)

	// Not found
	mock.ExpectQuery(query).
		WithArgs("missing", "reset_password", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err = repo.FindValidToken("missing", "reset_password")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// DB error
	mock.ExpectQuery(query).
		WithArgs("error", "reset_password", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("db error"))

	result, err = repo.FindValidToken("error", "reset_password")
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserTokenRepository_ConsumeToken(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()
	repo := repositorie.NewGormUserTokenRepository(db)

	userID := 1
	token := "abc123"
	purpose := "reset_password"

	// Success
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "user_tokens" SET "consumed"=$1,"used_at"=$2 WHERE user_id = $3 AND token = $4 AND purpose = $5 AND consumed = $6`)).
		WithArgs(true, sqlmock.AnyArg(), userID, token, purpose, false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.ConsumeToken(userID, token, purpose)
	assert.NoError(t, err)

	// Not found (RowsAffected=0)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "user_tokens" SET "consumed"=$1,"used_at"=$2 WHERE user_id = $3 AND token = $4 AND purpose = $5 AND consumed = $6`)).
		WithArgs(true, sqlmock.AnyArg(), userID, token, purpose, false).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.ConsumeToken(userID, token, purpose)
	assert.Error(t, err)
	assert.Equal(t, "token not found or already consumed", err.Error())

	// DB error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "user_tokens" SET "consumed"=$1,"used_at"=$2 WHERE user_id = $3 AND token = $4 AND purpose = $5 AND consumed = $6`)).
		WithArgs(true, sqlmock.AnyArg(), userID, token, purpose, false).
		WillReturnError(fmt.Errorf("db error"))
	mock.ExpectRollback()

	err = repo.ConsumeToken(userID, token, purpose)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserTokenRepository_DeleteExpired(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()
	repo := repositorie.NewGormUserTokenRepository(db)

	// Success
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_tokens" WHERE expires_at <= $1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteExpired()
	assert.NoError(t, err)

	// DB error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_tokens" WHERE expires_at <= $1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("delete error"))
	mock.ExpectRollback()

	err = repo.DeleteExpired()
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormUserTokenRepository_FindLatestValidToken(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()
	repo := repositorie.NewGormUserTokenRepository(db)

	userID := 1
	purpose := "reset_password"
	now := time.Now()
	token := &UserModel.UserToken{
		ID:        1,
		UserID:    userID,
		Token:     "abc123",
		Purpose:   purpose,
		Consumed:  false,
		ExpiresAt: now.Add(time.Hour),
	}

	// Include created_at and updated_at columns
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "token", "purpose", "consumed", "expires_at", "created_at", "updated_at",
	}).AddRow(token.ID, token.UserID, token.Token, token.Purpose, token.Consumed, token.ExpiresAt, now, now)

	// Success: allow the ORDER BY tiebreaker on id
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_tokens" WHERE user_id = $1 AND purpose = $2 AND consumed = $3 AND expires_at > $4 ORDER BY created_at DESC,"user_tokens"."id" LIMIT $5`,
	)).
		WithArgs(userID, purpose, false, sqlmock.AnyArg(), 1). // note: GORM passes LIMIT as a parameter
		WillReturnRows(rows)

	result, err := repo.FindLatestValidToken(userID, purpose)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token.Token, result.Token)

	// Not found
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_tokens" WHERE user_id = $1 AND purpose = $2 AND consumed = $3 AND expires_at > $4 ORDER BY created_at DESC,"user_tokens"."id" LIMIT $5`,
	)).
		WithArgs(userID, purpose, false, sqlmock.AnyArg(), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err = repo.FindLatestValidToken(userID, purpose)
	assert.NoError(t, err)
	assert.Nil(t, result)

	// DB error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_tokens" WHERE user_id = $1 AND purpose = $2 AND consumed = $3 AND expires_at > $4 ORDER BY created_at DESC,"user_tokens"."id" LIMIT $5`,
	)).
		WithArgs(userID, purpose, false, sqlmock.AnyArg(), 1).
		WillReturnError(fmt.Errorf("db error"))

	result, err = repo.FindLatestValidToken(userID, purpose)
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
