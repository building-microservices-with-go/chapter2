package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var programmingLanguages = map[string]string{
	"1": "C",
	"2": "C++",
	"3": "C#",
	"4": "Go",
	"5": "Java",
	"6": "Ruby",
}

type languageRequest struct {
	Language string
}

type languagesResponse struct {
	Languages map[string]string
}

type languageKey struct{}

func main() {
	http.HandleFunc("/languages", handleLanguages)
	http.HandleFunc("/languages/", handleLanguages)

	log.Println("Starting server on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleLanguages(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(rw, r)
		return
	case http.MethodPost:
		validateRequest(rw, r, handlePost)
		return
	case http.MethodPut:
		validateRequest(rw, r, handlePut)
		return
	default:
		http.NotFound(rw, r)
		return
	}
}

func handleGet(rw http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(rw)
	err := e.Encode(languagesResponse{
		Languages: programmingLanguages,
	})

	if err != nil {
		http.Error(rw, "Error encoding languages", http.StatusInternalServerError)
		return
	}
}

func validateRequest(rw http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	d := json.NewDecoder(r.Body)
	data := languageRequest{}

	err := d.Decode(&data)
	if err != nil {
		http.Error(rw, "Invalid request, invalid language request", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), languageKey{}, data.Language)
	req := r.WithContext(ctx)

	next(rw, req)
}

func handlePost(rw http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(languageKey{}).(string)

	id := strconv.Itoa(len(programmingLanguages) + 1)
	programmingLanguages[id] = data
}

func handlePut(rw http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(languageKey{}).(string)

	// parse the id
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	if id == "" {
		http.Error(rw, "Invalid request, invalid id", http.StatusBadRequest)
		return
	}

	programmingLanguages[id] = data
}
