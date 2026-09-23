package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kangy1103/chirpy-boot-dev/internal/auth"
	"github.com/Kangy1103/chirpy-boot-dev/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) CreateUsers(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("could not decode parameters: %s", err)
		ErrorResponse(w, 500, "Something went wrong")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("could not hash password: %s", err)
		ErrorResponse(w, 500, "Something went wrong")
		return
	}
	createUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), createUserParams)
	if err != nil {
		log.Printf("Could not create user: %s", err)
		ErrorResponse(w, 500, "Something went wrong")
		return
	}
	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	JSONResponse(w, 201, newUser)
	return
}
