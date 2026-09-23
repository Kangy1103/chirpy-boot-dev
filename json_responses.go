package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, code int, bodyJSON interface{}) {
	data, err := json.Marshal(bodyJSON)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func ErrorResponse(w http.ResponseWriter, code int, msg string) {
	type errorBody struct {
		Error string `json:"error"`
	}
	JSONResponse(w, code, errorBody{Error: msg})
}
