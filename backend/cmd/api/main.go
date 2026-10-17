package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Akinkfdr/cinema-ticketing/backend/internal/httpapi"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8000"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           httpapi.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("API dinleniyor: %s", addr)
	log.Fatal(server.ListenAndServe())
}
