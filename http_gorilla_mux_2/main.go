package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

type requestKey struct{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/languages", handleGet).Methods("GET")

	s := r.PathPrefix("/languages").Subrouter()
	s.Use(validationMiddleware)
	s.HandleFunc("", handlePost).Methods("POST")
	s.HandleFunc("/{id}", handlePut).Methods("PUT")

	http.Handle("/", r)
	log.Println("Starting server on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
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

func handlePost(rw http.ResponseWriter, r *http.Request) {
	req := r.Context().Value(requestKey{}).(languageRequest)

	id := strconv.Itoa(len(programmingLanguages) + 1)
	programmingLanguages[id] = req.Language
}

func handlePut(rw http.ResponseWriter, r *http.Request) {
	// parse the id
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	if id == "" {
		http.Error(rw, "Invalid request, invalid id", http.StatusBadRequest)
		return
	}

	req := r.Context().Value(requestKey{}).(languageRequest)

	programmingLanguages[id] = req.Language
}

func validationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		req := languageRequest{}

		err := d.Decode(&req)
		if err != nil {
			http.Error(rw, "Invalid request, invalid language request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), requestKey{}, req)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
