package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"calculator/backend/internal/api"
)

const defaultPort = "8080"
const defaultAllowedOrigin = "http://localhost:5173"

func main() {
	port := getenv("PORT", defaultPort)
	allowedOrigin := getenv("ALLOWED_ORIGIN", defaultAllowedOrigin)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	handler := api.CORSMiddleware(mux, allowedOrigin)

	addr := ":" + port
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("calculator backend listening on %s (allowed origin: %s)", addr, allowedOrigin)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
