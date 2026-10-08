package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/mahditd/url-shortener/bootstrap"
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

	flag.Parse()

	app := bootstrap.NewApp(*baseURL)

	server := &http.Server{
		Addr:         *addr,
		Handler:      app,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server.ListenAndServe()
}
