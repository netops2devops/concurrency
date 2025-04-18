package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/src-d/go-git.v4"
)

func RepoCloneWorker(id int, job <-chan string, result chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range job {
		log.Println("Worker ", id, " received ", url)
		repoName, found := strings.CutPrefix(url, "https://github.com/esnet/")
		if !found {
			result <- fmt.Sprintf("Worker %d failed trim prefix %s", id, url)
		}

		name, found := strings.CutSuffix(repoName, ".git")
		if !found {
			result <- fmt.Sprintf("Worker %d failed trim suffix %s", id, url)
		}

		clonepath := filepath.Join("/tmp", name)

		_, err := git.PlainClone(clonepath, false, &git.CloneOptions{
			URL: url,
		})
		if err != nil {
			result <- fmt.Sprintf("Worker %d failed to clone %s", id, url)
		}

		result <- fmt.Sprintf("Worker %d finished cloning %s", id, url)
	}
}
