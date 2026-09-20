package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	testServer()
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Received a request!")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
		hits)
	io.WriteString(w, body)
}

func (cfg *apiConfig) ResetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Metrics reset!")
}

type chirps struct {
	Body string `json:"body"`
}

func JSONDecoder(w http.ResponseWriter, req *http.Request) {
	type isValid struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(req.Body)
	chirp := chirps{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Could not decode chirp: %s", err)
		ErrorResponse(w, 500, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		ErrorResponse(w, 400, "Chirp is too long")
		return
	} else {
		isValid := isValid{
			Valid: true,
		}
		JSONResponse(w, 200, isValid)
		return
	}
}

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
