package main

import (
	"errors"
	"net/http"
)

func (conf *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if conf.platform != "dev" {
		err := errors.New("Operation allowed only in a dev env")
		respondWithError(w, http.StatusForbidden, "Operation allowed only in a dev env", err)
		return
	}
	conf.fileServerHits.Store(0)
	if err := conf.db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to reset DB", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
