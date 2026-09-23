package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Could not access Chirps: %s", err)
		ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	chirps := []ChirpResponse{}

	for _, chirp := range dbChirps {
		chirps = append(chirps, ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	JSONResponse(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) GetChirp(w http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Printf("Could not parse ID: %s", err)
		ErrorResponse(w, 404, "Something went wrong")
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirp(req.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(w, 404, "Something went wrong")
		return
	}
	if err != nil {
		log.Printf("Could not access Chirps: %s", err)
		ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	chirp := ChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	JSONResponse(w, http.StatusOK, chirp)
}
