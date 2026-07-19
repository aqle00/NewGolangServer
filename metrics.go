package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		fmt.Printf(" ++ incremented succesfully\n")
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// check if wrong method, only allow GET
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	w.WriteHeader(http.StatusOK)

	bodyText := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	var body = []byte(bodyText)
	w.Write(body)

	fmt.Printf("Hits: %v\n", cfg.fileserverHits.Load())
}
