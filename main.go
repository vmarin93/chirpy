package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vmarin93/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
}

func main() {
	godotenv.Load()
	const port = "8080"
	const filepathRoot = "."
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "prod"
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to the db")
	}
	dbQueries := database.New(db)
	conf := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secretKey:      secretKey,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/",
		conf.middlewareMetricsCount(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/chirps", conf.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", conf.handlerChirpsList)
	mux.HandleFunc("GET /api/chirps/{chirpID}", conf.handlerChirpsListOne)
	mux.HandleFunc("POST /api/users", conf.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", conf.handlerLogin)
	mux.HandleFunc("GET /admin/metrics", conf.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", conf.handlerReset)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
