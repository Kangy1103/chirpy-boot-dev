package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func testServer() {
	mux := NewServeMux()
	apiConfig := &apiConfig{}
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, http.StatusText(http.StatusOK))
	})

	mux.HandleFunc("POST /api/validate_chirp", JSONDecoder)

	mux.HandleFunc("GET /admin/metrics", apiConfig.Metrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.ResetMetrics)

	fmt.Println("Server started...")
	fs := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiConfig.MiddlewareMetricsInc(fs))
	fmt.Println("File server started...")
	log.Fatal(s.ListenAndServe())
}
