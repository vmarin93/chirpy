package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
