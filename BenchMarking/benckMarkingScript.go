package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	url        = "http://localhost:8080/load" // <-- change to your endpoint
	concurrent = 50
	duration   = 10 * time.Second
)

func main() {
	var wg sync.WaitGroup
	var total, success, fail int64
	var mu sync.Mutex

	deadline := time.Now().Add(duration)
	fmt.Println("Benchmarking", url, "for", duration, "with", concurrent, "concurrent users...")

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{}

			for time.Now().Before(deadline) {
				resp, err := client.Get(url)

				mu.Lock()
				total++
				if err != nil || resp.StatusCode != http.StatusOK {
					fail++
				} else {
					io.Copy(io.Discard, resp.Body) // discard response body
					resp.Body.Close()
					success++
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Println("Done.")
	fmt.Printf("Total Requests: %d\n", total)
	fmt.Printf("Successful:     %d\n", success)
	fmt.Printf("Failed:         %d\n", fail)
	fmt.Printf("Success Rate:   %.2f%%\n", float64(success)/float64(total)*100)
}
