package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var workerFlag = flag.Int("workers", 4, "number of concurrent workers available")
	flag.Parse()
	if *workerFlag < 1 {
		return errors.New("need greater than 0 number of workers")
	}
	if *workerFlag > 20 {
		return errors.New("max 20 workers allowed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urlList, err := GetURLList(flag.Arg(0))
	if err != nil {
		//file invalid somehow
		return err
	}

	for _, url := range urlList {
		if !isValidURL(url) {
			continue
		}
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			url,
			nil,
		)
		if err != nil {
			// invalid request / URL
			fmt.Println("Invalid request")
			return err
		}

		client := &http.Client{}

		timePreCall := time.Now()
		resp, err := client.Do(req)
		travelTime := time.Since(timePreCall)
		if err != nil {
			// request failed or context timed out
			fmt.Println("Request failed/ timed out")
			return err
		}

		fmt.Printf("%s         | Status Code: %d | Travel time: %v\n", url, resp.StatusCode, travelTime.Round(time.Millisecond))

		defer resp.Body.Close()
	}

	return nil
}
