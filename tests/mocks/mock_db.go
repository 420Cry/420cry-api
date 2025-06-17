package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB creates a mocked gorm.DB with sqlmock and returns a cleanup function.
func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
	}

	return gormDB, mock, cleanup
}
