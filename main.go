package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	var apiCfg = apiConfig{
		fileserverHits: atomic.Int32{},
	}

	//a ServeMux is basically just a handler that maps a url endpoint to a handler function
	// it does so by matching the closest pattern, instead of exact pattern
	// make new servemux here
	mux := http.NewServeMux()

	//register endpoint handler to serve files or functions
	// hold ctrl and hover over mux.Handle(), .Handler() functions to see doc and read how to use
	//StripPrefix() srtips /app, so the pattern is removed from the file path returned from calling .FileServer()
	// but the pattern on the web is still /app/, so the user has to use that endpoing instead of /
	mux.Handle("/app/", apiCfg.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.HandleFunc("GET /healthz", handleReadiness)
	mux.HandleFunc("GET /metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.HandleFunc("POST /reset", http.HandlerFunc(apiCfg.handlerResetHits))

	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ==================================

	//make server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// print stuff
	log.Printf("Serving on port: %s\n", port)

	// here, srv.ListenAndServe() is called, starts running and block forever, until the server closes
	// if the server is closed either:
	// we closed it, we want to log this just to keep track
	// something wrong happened, need to track for sure
	// that is the reason log.Fatal() is used to wrap the function, because we don't really ever want to close the server, and if it does we want to know the reason why
	log.Fatal(srv.ListenAndServe())
}
