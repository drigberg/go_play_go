package main

import (
	"log"
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
	log.Println("I exist! (port: " + port + ", host: " + host + ")")
}
