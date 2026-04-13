package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vmarin93/chirpy/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	if err != nil {
		log.Fatal("Unable to connect to the db")
	}
	const port = "8080"
	const filepathRoot = "."
	conf := &apiConfig{db: dbQueries}
	mux := http.NewServeMux()
	mux.Handle("/app/",
		conf.middlewareMetricsCount(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidation)
	mux.HandleFunc("GET /admin/metrics", conf.getMetricsCount)
	mux.HandleFunc("POST /admin/reset", conf.resetMetricsCount)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
