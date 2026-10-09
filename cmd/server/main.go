package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mahditd/url-shortener/bootstrap"
	"github.com/mahditd/url-shortener/internal/domain/ports"
	"github.com/mahditd/url-shortener/internal/infrastructure/database"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/postgres"
	"github.com/mahditd/url-shortener/internal/presentation/middleware"
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

	rateLimiter := middleware.NewRateLimiter(
		10,
		time.Minute,
	)

	app := bootstrap.NewApp(*baseURL, repository , rateLimiter)

	server := &http.Server{
		Addr:         *addr,
		Handler:      app,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		panic(err)
	}
}
