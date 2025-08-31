package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Concurrent Web Fetcher with Timeout

// Write a Go program that:
// - Takes a list of URLs as CLI arguments
// - Spawns a goroutine for each URL to fetch it concurrently.
// - Uses a channel to collect the results (either success or error).
// - Uses a context with timeout so that if any fetch takes too long, it gets canceled automatically.
// - Prints out all results before exiting.

type Result struct {
	Url    string
	Status string
	Err    error
}

func FetchURL(ctx context.Context, url string, ch chan<- Result) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- Result{Err: err, Url: url}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- Result{Err: err, Url: url}
		return
	}
	defer resp.Body.Close()

	ch <- Result{Url: url, Status: resp.Status, Err: nil}
}

func main() {
	ch := make(chan Result)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, url := range os.Args[1:] {
		go FetchURL(ctx, url, ch)
	}

	for range os.Args[1:] {
		result := <-ch
		if result.Err != nil {
			fmt.Println(result.Url, result.Err)
		}
		if result.Err == nil {
			fmt.Println(result.Url, result.Status)
		}
	}
}
