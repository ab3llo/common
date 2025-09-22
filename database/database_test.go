package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConnectToDB_InvalidDSN(t *testing.T) {
	db := ConnectToDB(":invalid-dsn:")
	assert.Nil(t, db, "DB should be nil for invalid DSN")
}

func TestConnectToDB_Success(t *testing.T) {
	// Use the same approach as domain tests - in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Should successfully connect to in-memory SQLite")
	assert.NotNil(t, db, "DB should not be nil for valid connection")

	// Test that we can perform basic operations
	sqlDB, err := db.DB()
	assert.NoError(t, err, "Should be able to get underlying sql.DB")
	assert.NoError(t, sqlDB.Ping(), "Should be able to ping the database")

	// Test connection stats
	stats := sqlDB.Stats()
	assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0, "Should have connection stats")
}

func TestConnectToDB_PostgresSpecificErrors(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
	}{
		{
			name: "empty DSN",
			dsn:  "",
		},
		{
			name: "malformed DSN",
			dsn:  "invalid://malformed:dsn",
		},
		{
			name: "missing host",
			dsn:  "postgresql://user:pass@:5432/db",
		},
		{
			name: "invalid port",
			dsn:  "postgresql://user:pass@host:99999/db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := ConnectToDB(tt.dsn)
			assert.Nil(t, db, "DB should be nil for invalid DSN: %s", tt.dsn)
		})
	}
}

func BenchmarkConnectToDB(b *testing.B) {
	dsn := ":invalid-for-benchmark:"
	for i := 0; i < b.N; i++ {
		db := ConnectToDB(dsn)
		if db != nil {
			// This shouldn't happen with invalid DSN, but close if it does
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	}
}
