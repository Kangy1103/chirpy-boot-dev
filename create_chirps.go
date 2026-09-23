package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kangy1103/chirpy-boot-dev/internal/database"
	"github.com/google/uuid"
)

type ChirpRequest struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) ChirpDecoder(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := ChirpRequest{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Could not decode chirp: %s", err)
		ErrorResponse(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		ErrorResponse(w, 400, "Chirp is too long")
		return
	}
	splitChirp := strings.Split(chirp.Body, " ")
	cleanChirp := []string{}
	for _, word := range splitChirp {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			cleanChirp = append(cleanChirp, "****")
		default:
			cleanChirp = append(cleanChirp, word)
		}
	}
	cleanedBody := strings.Join(cleanChirp, " ")
	userID, err := uuid.Parse(chirp.UserID)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid user_id")
		return
	}
	newChirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	}
	dbChirp, err := cfg.dbQueries.CreateChirp(req.Context(), newChirpParams)
	if err != nil {
		log.Printf("Could not create chirp: %s", err)
		ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	chirpResponse := ChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	JSONResponse(w, 201, chirpResponse)
	return
}
