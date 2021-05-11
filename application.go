package main

import (
	"os"
)

func parseEnv() (string, string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	return port, host
}

func main() {
	port, host := parseEnv()
	RunServer(host, port)
}
