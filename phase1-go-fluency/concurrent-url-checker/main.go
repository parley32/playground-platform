package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/tabwriter"
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

	urlList, err := getURLList(flag.Arg(0))
	if err != nil {
		fmt.Println("file invalid somehow")
		return err
	}

	jobs := make(chan Job)
	results := make(chan Result)
	client := &http.Client{}

	var wg sync.WaitGroup
	for range *workerFlag {
		wg.Add(1)
		go startWorker(ctx, jobs, results, client, &wg)
	}

	go func() {
		for _, url := range urlList {
			jobs <- Job{URL: normalizeURL(url)}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "URL\tSTATUS\tDURATION\tERROR")
	fmt.Fprintln(w, "-----------------------------------\t\t\t")

	for result := range results {
		if result.Err != nil {
			fmt.Fprintf(w, "%s\t-\t%s\t%v\n",
				result.URL,
				result.Duration,
				result.Err,
			)
			continue
		}

		fmt.Fprintf(w, "%s\t%d\t%s\t\n",
			result.URL,
			result.StatusCode,
			result.Duration,
		)
	}
	w.Flush()

	return nil
}
