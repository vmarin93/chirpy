package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/vmarin93/chirpy/internal/auth"
)

func (conf *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		AccessToken string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse refresh token", err)
		return
	}
	refreshTokenDB, err := conf.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not found in db", err)
		return
	}
	if refreshTokenDB.RevokedAt.Valid {
		err = errors.New("Refresh Token revoked")
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", err)
		return
	}
	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		err = errors.New("Refresh Token has expired")
		respondWithError(w, http.StatusUnauthorized, "Refresh Token has expired", err)
		return
	}
	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(refreshTokenDB.UserID, conf.secretKey, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create accessToken", err)
		return
	}
	respondWithJson(w, http.StatusOK, responseBody{AccessToken: accessToken})
}

func (conf *apiConfig) handlerRefreshTokenRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse refresh token", err)
		return
	}
	if err := conf.db.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
