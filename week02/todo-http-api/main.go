package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, output io.Writer) error {
	mux := &http.ServeMux{}

	taskHandler := NewTaskHandler(NewTaskService())
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
