package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vmarin93/chirpy/internal/auth"
	"github.com/vmarin93/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (conf *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthenticated users can't post chirps", err)
		return
	}
	userID, err := auth.ValidateJWT(token, conf.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthenticated users can't post chirps", err)
		return
	}
	validBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	chirp, err := conf.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   validBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to add chirp to DB", err)
		return
	}
	respondWithJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (conf *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}
	sortDirection, err := sortDirectionFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid sort direction", err)
		return
	}
	if authorID != uuid.Nil {
		chirps, err = conf.db.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		chirps, err = conf.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps from DB", err)
		return
	}
	responseChirps := []Chirp{}
	for _, chirp := range chirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	if sortDirection == "desc" {
		slices.SortFunc(responseChirps, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}
	respondWithJson(w, http.StatusOK, responseChirps)
}

func (conf *apiConfig) handlerChirpsListOne(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse URL param", err)
		return
	}
	chirp, err := conf.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found in the DB", err)
		return
	}
	respondWithJson(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (conf *apiConfig) handlerChirpsDeleteOne(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT token not found in header", err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, conf.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse URL param", err)
		return
	}
	chirp, err := conf.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found in the DB", err)
		return
	}
	if chirp.UserID != userID {
		err = errors.New("You can only delete your own chirps")
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps", err)
		return
	}
	if err := conf.db.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete chirp from DB", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateChirp(body string) (string, error) {
	const maxChirpLen = 140
	if len(body) > maxChirpLen {
		return "", errors.New("Chirp is too long")
	}
	sanitized := sanitize(body)
	return sanitized, nil
}

func sanitize(s string) string {
	censored := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(s, " ")
	for i, word := range words {
		if _, ok := censored[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func sortDirectionFromRequest(r *http.Request) (string, error) {
	sortDirection := r.URL.Query().Get("sort")
	if sortDirection == "" {
		return "", nil
	}
	if sortDirection != "asc" && sortDirection != "desc" {
		return "", errors.New("Bad sorting parameter")
	}
	return sortDirection, nil
}
