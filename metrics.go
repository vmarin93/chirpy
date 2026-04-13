package main

import (
	"fmt"
	"net/http"
)

func (conf *apiConfig) middlewareMetricsCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (conf *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf([]byte{},
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, conf.fileServerHits.Load()))
}
