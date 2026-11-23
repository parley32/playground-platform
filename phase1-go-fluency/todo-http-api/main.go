package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	mux := &http.ServeMux{}

	taskHandler := NewTaskHandler(NewTaskService(fileStore{}))
	mux.HandleFunc("POST /task/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/task/" {
			http.NotFound(w, r)
			return
		}
		taskHandler.Create(w, r)
	})

	mux.HandleFunc("GET /task/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/task/" {
			http.NotFound(w, r)
			return
		}
		taskHandler.List(w)
	})

	mux.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) {
		taskHandler.ListSingleItem(w, r.PathValue("id"))
	})

	mux.HandleFunc("DELETE /task/{id}", func(w http.ResponseWriter, r *http.Request) {
		taskHandler.Delete(w, r.PathValue("id"))
	})

	mux.HandleFunc("PATCH /task/{id}", func(w http.ResponseWriter, r *http.Request) {
		taskHandler.Check(w, r.PathValue("id"))
	})

	return http.ListenAndServe(":8080", mux)
}
