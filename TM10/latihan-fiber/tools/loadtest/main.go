package main

import (
	"sync"
)

func fire(n int, job func(i int) string) map[string]int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int, out chan<- string) {
			defer wg.Done()

			<-start
			out <- job(i)
		}(i, results)
	}

	close(start) // buka gerbang
	wg.Wait()    // tunggu SEMUA goroutine selesai
	close(results)

	tally := map[string]int{}

	for r := range results {
		tally[r]++
	}

	return tally
}

func main() {
	results := fire(100, func(i int) string {
		if i%2 == 0 {
			return "success"
		}
		return "failed"
	})

	for status, count := range results {
		println(status, count)
	}
}
