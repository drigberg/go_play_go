package main

import (
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	RunServer(port)
}
