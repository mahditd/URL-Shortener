package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/mahditd/url-shortener/bootstrap"
	"github.com/mahditd/url-shortener/internal/domain/ports"
	"github.com/mahditd/url-shortener/internal/infrastructure/database"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/postgres"
)

func main() {

	addr := flag.String(
		"addr",
		":8080",
		"server listen address",
	)

	baseURL := flag.String(
		"base",
		"http://localhost:8080",
		"base URL for generated short links",
	)

	storage := flag.String(
		"storage",
		"memory",
		"storage backend: memory or postgres",
	)

	flag.Parse()

	var repository ports.LinkRepository

	switch *storage {
	case "memory":

		repository = memory.NewLinkRepository()

	case "postgres":
		db, err := database.NewPostgresConnection()

		if err != nil {
			panic(err)
		}
		err = database.AutoMigrate(db)

		if err != nil {
			panic(err)
		}

		repository = postgres.NewLinkRepository(db)
	default:
		panic("invalid storage option")
	}

	app := bootstrap.NewApp(*baseURL, repository)

	server := &http.Server{
		Addr:         *addr,
		Handler:      app,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server.ListenAndServe()
}
