package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) ResetMetrics(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		ErrorResponse(w, 403, "Forbidden")
		return
	}
	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		log.Printf("Could not access database: %s", err)
		ErrorResponse(w, 500, "Something went wrong")
		return
	}
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Metrics reset!")
}
