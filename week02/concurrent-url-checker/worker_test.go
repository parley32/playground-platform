package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDoRequestValidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := doRequest(ctx, server.Client(), server.URL)

	if result.Err != nil {
		t.Fatalf("doRequest returned error for valid URL: %v", result.Err)
	}
	if result.URL != server.URL {
		t.Fatalf("expected URL %q, got %q", server.URL, result.URL)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, result.StatusCode)
	}
	if result.Duration <= 0 {
		t.Fatal("expected a positive duration")
	}
}

func TestDoRequestReturns404Result(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := doRequest(ctx, server.Client(), server.URL)

	if result.Err != nil {
		t.Fatalf("expected 404 to be handled as a valid response, got error: %v", result.Err)
	}
	if result.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, result.StatusCode)
	}
	if result.URL != server.URL {
		t.Fatalf("expected URL %q, got %q", server.URL, result.URL)
	}
}

func TestDoRequestInvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := doRequest(ctx, &http.Client{}, "://bad-url")

	if result.Err == nil {
		t.Fatal("expected invalid URL to return an error")
	}
	if result.StatusCode != 0 {
		t.Fatalf("expected zero status code for invalid request, got %d", result.StatusCode)
	}
}
