package database

import "testing"

func TestAutoMigrate(t *testing.T) {
	db, err := NewPostgresConnection()

	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	err = AutoMigrate(db)

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}
}