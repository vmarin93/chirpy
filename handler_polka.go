package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/vmarin93/chirpy/internal/auth"
)

func (conf *apiConfig) handlerPolkaUpgradeHook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetPolkaAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad authorization in header", err)
		return
	}
	if apiKey != conf.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Bad API key", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	user, err := conf.db.UpgradeUserMembership(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User doesn't exist", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Unable to upgrade user membership", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Membership upgraded for user: %v", user.ID)
}
