package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// clean json body
	cleanedBody := badWordFilter(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_body: cleanedBody,
	})
}

func badWordFilter(s string) string {
	// map of empty strings that contains bad words as keys

	badwords := map[string]struct{}{
		//register bad words
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	//split string into words []string
	words := strings.Split(s, " ")

	// loop through wordsMap to check if an input word is a bad word, replace with **** if so
	for i, word := range words {
		if _, exist := badwords[strings.ToLower(word)]; exist {
			words[i] = "****"
		}
	}
	// after loop, the words []string variable now contains filtered versions of bad words at the indexes of replaced bad words
	// join the []string into 1 string, with a space " " as seperator between the words
	cleanedString := strings.Join(words, " ")
	return cleanedString
}
