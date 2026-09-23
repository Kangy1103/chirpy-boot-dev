package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Kangy1103/chirpy-boot-dev/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file not found: %s", err)
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("cannot load database: %s", err)
	}
	dbQueries := database.New(db)
	cfg := &apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", isServerReady)

	mux.HandleFunc("POST /api/users", cfg.CreateUsers)
	mux.HandleFunc("POST /api/chirps", cfg.ChirpDecoder)
	mux.HandleFunc("POST /api/login", cfg.UserLogin)
	mux.HandleFunc("POST /admin/reset", cfg.ResetMetrics)

	mux.HandleFunc("GET /admin/metrics", cfg.Metrics)
	mux.HandleFunc("GET /api/chirps", cfg.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirp)

	fmt.Println("Server started...")
	fs := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(fs))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("File server started...")
	log.Fatal(s.ListenAndServe())
}
