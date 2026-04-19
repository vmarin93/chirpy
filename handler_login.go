package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vmarin93/chirpy/internal/auth"
	"github.com/vmarin93/chirpy/internal/database"
)

func (conf *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		User
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params", err)
		return
	}
	user, err := conf.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, conf.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create JWT token", err)
		return
	}
	refreshToken, err := conf.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to add refresh token to DB", err)
		return
	}
	respondWithJson(w, http.StatusOK, responseBody{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}
