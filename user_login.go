package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Kangy1103/chirpy-boot-dev/internal/auth"
)

func (cfg *apiConfig) UserLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("could not decode paramaters: %s", err)
		ErrorResponse(w, 500, "Something went wrong")
		return
	}
	dbUser, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(w, 401, "Incorrect email or password")
		return
	}
	if err != nil {
		log.Printf("Could not access user: %s", err)
		ErrorResponse(w, 401, "Incorrect email or password")
		return
	}
	password, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		ErrorResponse(w, 401, "Incorrect email or password")
		return
	}
	if !password {
		ErrorResponse(w, 401, "Incorrect email or password")
		return
	}
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	JSONResponse(w, http.StatusOK, user)
}
