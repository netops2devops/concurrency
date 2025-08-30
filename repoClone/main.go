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
	MAX_WORKERS = 5
)

func CollectPublicRepoURLs(org string) []string {
	var c http.Client
	resp, err := c.Get(fmt.Sprintf("https://api.github.com/orgs/%s/repos", org))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var result []map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal("Unable to unmarshal response", err)
	}

	urls := make([]string, 0, 100)

	// Extract repo url
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

	// Create workerpool to clone repo
	for i := 1; i <= MAX_WORKERS; i++ {
		wg.Add(1)
		go RepoCloneWorker(i, jobs, results, &wg)
	}

	// Send all urls to job queue channel
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// Collect results as they become available
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r)
	}

	// Cleanup all cloned repositories
	RepoCleaner()
}

func RepoCleaner() bool {
	err := os.RemoveAll("/tmp/demo")
	if err != nil {
		log.Fatal(err)
	}
	return true
}
