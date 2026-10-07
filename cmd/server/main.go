package main

import (
	"flag"

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

	app.Run(*addr)
}
