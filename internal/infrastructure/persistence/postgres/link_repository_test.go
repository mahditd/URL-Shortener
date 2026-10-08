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
