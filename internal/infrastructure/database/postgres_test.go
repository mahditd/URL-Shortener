package database

import (
	"testing"
)

func TestPostgresConnection(t *testing.T) {
	db, err := NewPostgresConnection()

	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("database ping failed: %v", err)
	}
}
