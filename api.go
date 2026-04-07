package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {
	const maxChirpLen = 140
	censored := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Sanitized string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode json in request body", err)
	}
	if len(params.Body) > maxChirpLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	} else {
		respondWithJson(w, http.StatusOK, returnVals{
			Sanitized: sanitize(params.Body, censored)})
	}
}

func sanitize(s string, censored map[string]struct{}) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		if _, ok := censored[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON when responding %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with a 5xx error %s", msg)
	}
	respondWithJson(w, code, errorResponse{Error: msg})
}
