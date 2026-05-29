package main

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ezenkico/nats-auth-redirect/redirects"
)

func main() {
	godotenv.Load()
	server := os.Getenv("AUTH_SERVER")
	if server != "" {
		println("Setting up http target")
		redirects.SetupHttpAuth(server)
		return
	}
	println("No listener specified")
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
