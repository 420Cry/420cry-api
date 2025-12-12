// Package suites provides test suite functionality for the application.
package suites

import (
	"os"
	"testing"

	"cry-api/app/config"
	"cry-api/app/container"
	"cry-api/app/logger"
	"cry-api/app/models"
	EnvTypes "cry-api/app/types/env"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserTestSuite provides a test suite for user-related functionality
type UserTestSuite struct {
	suite.Suite
	db        *gorm.DB
	container *container.Container
	logger    *logger.Logger
}

// SetupSuite initializes the test suite
func (suite *UserTestSuite) SetupSuite() {
	// Set test environment variables first (JWT_SECRET must be set before any JWT package imports)
	_ = os.Setenv("APP_ENV", "test")
	_ = os.Setenv("JWT_SECRET", "testsecretkey123456789012345678901234567890") // Set before JWT package init
	_ = os.Setenv("LOG_LEVEL", "error")                                        // Reduce log noise during tests
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "3306")
	_ = os.Setenv("DB_DATABASE", "test_db")
	_ = os.Setenv("DB_USERNAME", "test_user")
	_ = os.Setenv("DB_PASSWORD", "test_password")
	_ = os.Setenv("NO_REPLY_EMAIL", "noreply@test.com")
	_ = os.Setenv("CRY_APP_URL", "http://localhost:3000")
	_ = os.Setenv("CRY_API_URL", "http://localhost:8080")
	_ = os.Setenv("SMTP_HOST", "localhost")
	_ = os.Setenv("SMTP_PORT", "1025")
	_ = os.Setenv("API_PORT", "8080")

	// Set test configuration directly
	testCfg := &EnvTypes.EnvConfig{
		AppEnv:     "test",
		CryAppURL:  "http://localhost:3000",
		CryAPIURL:  "http://localhost:8080",
		APIPort:    8080,
		DBHost:     "localhost",
		DBPort:     3306,
		DBDatabase: "test_db",
		DBUserName: "test_user",
		DBPassword: "test_password",
		SMTPConfig: EnvTypes.SMTPConfig{
			Host: "localhost",
			Port: "1025",
		},
		NoReplyEmail: "noreply@test.com",
	}
	config.SetTestConfig(testCfg)

	// Initialize logger
	suite.logger = logger.GetLogger()

	// Setup in-memory SQLite database for tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.UserToken{})
	suite.Require().NoError(err)

	// Initialize container with test dependencies
	suite.container = container.InitializeContainer(testCfg, db)

	suite.logger.Info("Test suite setup completed")
}

// TearDownSuite cleans up after the test suite
func (suite *UserTestSuite) TearDownSuite() {
	// Close database connection
	if sqlDB, err := suite.db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	suite.logger.Info("Test suite teardown completed")
}

// SetupTest runs before each test
func (suite *UserTestSuite) SetupTest() {
	// Clean up database before each test
	suite.db.Exec("DELETE FROM user_tokens")
	suite.db.Exec("DELETE FROM users")
}

// TearDownTest runs after each test
func (suite *UserTestSuite) TearDownTest() {
	// Clean up database after each test
	suite.db.Exec("DELETE FROM user_tokens")
	suite.db.Exec("DELETE FROM users")
}

// GetContainer returns the test container
func (suite *UserTestSuite) GetContainer() *container.Container {
	return suite.container
}

// GetDB returns the test database
func (suite *UserTestSuite) GetDB() *gorm.DB {
	return suite.db
}

// GetLogger returns the test logger
func (suite *UserTestSuite) GetLogger() *logger.Logger {
	return suite.logger
}

// RunUserTestSuite runs the user test suite
func RunUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
