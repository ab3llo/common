package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToDB_InvalidDSN(t *testing.T) {
	db := ConnectToDB(":invalid-dsn:")
	assert.Nil(t, db, "DB should be nil for invalid DSN")
}

// Note: Testing a successful DB connection would require a real or mock database.
// For unit tests, you can use a test Postgres instance or mock gorm.Open if needed.
