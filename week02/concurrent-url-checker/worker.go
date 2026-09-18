package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func startWorker(ctx context.Context, jobs <-chan Job, results chan<- Result, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			result := doRequest(ctx, client, job.URL)

			select {
			case <-ctx.Done():
				return
			case results <- result:
			}
		}
	}
}

func prepareRequestBody(ctx context.Context, url string) (*http.Request, error) {
	return http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
}

func doRequest(ctx context.Context, client *http.Client, url string) Result {
	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if !isValidURL(url) {
		return Result{URL: url, Err: fmt.Errorf("Invalid url")}
	}

	req, err := prepareRequestBody(requestCtx, url)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	timePreCall := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(timePreCall)
	if err != nil {
		// request failed or context timed out
		return Result{URL: url, Err: fmt.Errorf("HTTP Request failed: %w", err)}
	}
	defer resp.Body.Close()
	result := Result{URL: string(url),
		StatusCode: resp.StatusCode,
		Duration:   duration.Round(time.Millisecond),
		Err:        nil}
	return result
}
