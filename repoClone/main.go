package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	MAX_WORKERS = 10
)

func CollectPublicRepoURLs(org string) []string {
	var c http.Client
	resp, err := c.Get(fmt.Sprintf("https://api.github.com/orgs/%s/repos", org))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// read client response
	body, err := io.ReadAll(resp.Body)

	var result []map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal("Unable to unmarshal response", err)
	}

	urls := make([]string, 0, 100) // length=0, capacity=100

	// extract repo url
	for _, repo := range result {
		u, ok := repo["clone_url"].(string)
		if !ok {
			continue
		}
		urls = append(urls, u)
	}
	return urls
}

func main() {
	org := os.Getenv("ORG")
	urls := CollectPublicRepoURLs(org)

	results := make(chan string, len(urls))
	jobs := make(chan string, len(urls))
	var wg sync.WaitGroup

	// Start workerpool to clone repo
	for i := 1; i <= MAX_WORKERS; i++ {
		wg.Add(1)
		go RepoCloneWorker(i, jobs, results, &wg)
	}

	// Send all urls to job queue
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)
	wg.Wait()

	// Collect results
	close(results)
	for r := range results {
		fmt.Println(r)
	}
}
