package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tbone317/aa-secretary/internal/db"
)

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "app.db"
	}

	sqlDB, err := db.Open(dsn)
	if err != nil {
		log.Fatalf("database open failed: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("database close failed: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Migrate(ctx, sqlDB, "internal/db/migrations"); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	fmt.Printf("Starting server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
