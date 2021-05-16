package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func ServeClientApp() {
	r := http.NewServeMux()
	buildHandler := http.FileServer(http.Dir("app/build"))
	r.Handle("/", buildHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("App server started on PORT 8080")
	log.Fatal(srv.ListenAndServe())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if os.Getenv("ENV") == "PRODUCTION" {
		go ServeClientApp()
	}
	RunWebsocketServer()
}
