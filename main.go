package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/aqle00/NewGolangServer.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	platform       string
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	var apiCfg = apiConfig{
		platform:       os.Getenv("PLATFORM"),
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
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
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(apiCfg.handlerReset))
	mux.HandleFunc("POST /admin/resethits", http.HandlerFunc(apiCfg.handlerResetHits))
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

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
