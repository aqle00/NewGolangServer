package main

import "net/http"

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	//set header content-type
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if r.Method != "GET" {
		w.WriteHeader(405)
	}
	// set response status
	w.WriteHeader(200)
	// set response message "OK"
	var body = []byte(http.StatusText(http.StatusOK))
	w.Write(body)
}
