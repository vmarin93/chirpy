package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	conf := &apiConfig{}
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
