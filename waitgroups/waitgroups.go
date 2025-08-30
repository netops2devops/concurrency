package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	ch := make(chan string, 5)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Go(func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			resp, err := MakeHTTPRequest(ctx, "http://httpbin.org/status/404,200,401,301,500,200,400,403")
			if err != nil {
				log.Println(err)
				cancel()
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode == http.StatusInternalServerError {
				log.Println("Internal Server Error → cancel just this request")
				cancel() // cancels only this goroutine’s context
				return
			}
			ch <- resp.Status
		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i)
	}
}

func MakeHTTPRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
