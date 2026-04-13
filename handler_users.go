package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (conf *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type responseBody struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params", err)
	}
	user, err := conf.db.CreateUsers(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
	}
	respondWithJson(w, http.StatusCreated, responseBody{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
