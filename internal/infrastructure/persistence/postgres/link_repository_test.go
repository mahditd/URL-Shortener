package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/mahditd/url-shortener/internal/domain/entities"
	"github.com/mahditd/url-shortener/internal/infrastructure/database"
)

func TestLinkRepositorySaveAndFind(t *testing.T) {
	db, err := database.NewPostgresConnection()

	if err != nil {
		t.Fatalf("database connection failed: %v", err)
	}

	err = database.AutoMigrate(db)

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	repo := NewLinkRepository(db)

	sqlDB, err := db.DB()

	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	defer func() {
		sqlDB.Exec("DELETE FROM links")
	}()

	link := entities.Link{
		Code:      fmt.Sprintf("test_%d", time.Now().UnixNano()),
		URL:       "https://example.com",
		CreatedAt: time.Now(),
	}

	err = repo.Save(link)

	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	found, err := repo.FindByCode(link.Code)

	if err != nil {
		t.Fatalf("find failed: %v", err)
	}

	if found.URL != link.URL {
		t.Fatalf("expected %s got %s", link.URL, found.URL)
	}
}

func TestLinkRepositoryRestartSimulation(t *testing.T) {
	db, err := database.NewPostgresConnection()

	if err != nil {
		t.Fatalf("database connection failed: %v", err)
	}

	err = database.AutoMigrate(db)

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	repo := NewLinkRepository(db)

	sqlDB, err := db.DB()

	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	defer sqlDB.Close()

	defer func() {
		sqlDB.Exec("DELETE FROM links")
	}()

	link := entities.Link{
		Code:      fmt.Sprintf("restart_%d", time.Now().UnixNano()),
		URL:       "https://restart-test.com",
		CreatedAt: time.Now(),
	}

	err = repo.Save(link)

	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	db2, err := database.NewPostgresConnection()

	if err != nil {
		t.Fatalf("database connection failed: %v", err)
	}

	err = database.AutoMigrate(db2)

	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	repo2 := NewLinkRepository(db2)

	sqlDB2, err := db2.DB()

	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	defer sqlDB2.Close()

	defer func() {
		sqlDB2.Exec("DELETE FROM links")
	}()

	found, err := repo2.FindByCode(link.Code)

	if err != nil {
		t.Fatalf("find after restart failed: %v", err)
	}

	if found.URL != link.URL {
		t.Fatalf("expected %s got %s", link.URL, found.URL)
	}

}
