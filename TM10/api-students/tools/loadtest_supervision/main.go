package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func fire(n int, job func(i int) (int, string)) map[string]int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()
			<-start

			status, result := job(idx)
			results <- fmt.Sprintf("%d %s", status, result)
		}(i)
	}

	close(start)
	wg.Wait()
	close(results)

	tally := make(map[string]int)
	for result := range results {
		tally[result]++
	}

	return tally
}

func main() {
	supervisorID := 1
	totalStudents := 30

	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		fmt.Println("ADMIN_TOKEN belum diatur.")
		fmt.Println(`Jalankan: $env:ADMIN_TOKEN="TOKEN_JWT_ADMIN"`)
		return
	}

	fmt.Printf(
		"=== UJI BEBAN D.3: %d Mahasiswa Berebut Pembimbing ID %d ===\n",
		totalStudents,
		supervisorID,
	)

	startTime := time.Now()

	results := fire(totalStudents, func(i int) (int, string) {
		studentID := i + 6

		payload, err := json.Marshal(map[string]int{
			"student_id": studentID,
		})
		if err != nil {
			return 500, "JSON_ERROR"
		}

		url := fmt.Sprintf(
			"http://localhost:3000/api/supervisors/%d/assign",
			supervisorID,
		)

		req, err := http.NewRequest(
			http.MethodPost,
			url,
			bytes.NewBuffer(payload),
		)
		if err != nil {
			return 500, "REQUEST_ERROR"
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return 500, "CONNECTION_ERROR"
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		responseBody := string(body)

		switch resp.StatusCode {
		case http.StatusCreated:
			return 201, "CREATED"
		case http.StatusConflict:
			return 409, "CONFLICT_OR_FULL"
		case http.StatusUnauthorized:
			return 401, "UNAUTHORIZED"
		case http.StatusForbidden:
			return 403, "FORBIDDEN"
		default:
			return resp.StatusCode, responseBody
		}
	})

	duration := time.Since(startTime)

	fmt.Println("\nHasil Uji Beban:")
	for result, count := range results {
		fmt.Printf("  %-30s : %d\n", result, count)
	}

	fmt.Printf("Waktu total: %v\n", duration)
}
