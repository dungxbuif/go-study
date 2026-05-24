package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func FetchData(ctx context.Context, client *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FetchMultiple(urls []string, timeout time.Duration) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []string
	var fetchErr error

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			data, err := FetchData(ctx, client, u)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fetchErr = err
			} else {
				results = append(results, data)
			}
		}(url)
	}

	wg.Wait()
	if fetchErr != nil {
		return nil, fetchErr
	}
	return results, nil
}

func main() {
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/status/200",
	}

	fmt.Println("=== Fetching with 3s Timeout ===")
	results, err := FetchMultiple(urls, 3*time.Second)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	} else {
		fmt.Printf("Successfully fetched %d responses\n", len(results))
	}

	fmt.Println("\n=== Fetching with 500ms Timeout (expect fail) ===")
	_, err = FetchMultiple(urls, 500*time.Millisecond)
	if err != nil {
		fmt.Printf("Expected error occurred: %v\n", err)
	}
}
