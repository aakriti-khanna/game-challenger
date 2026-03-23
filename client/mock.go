package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"game-challenge/models"
	"game-challenge/utils"
)

// MockUserEngine stimulates thousands of users answering and hitting the API concurrently
func MockUserEngine(numUsers int, targetURL string) {
	var wg sync.WaitGroup

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     numUsers,
			MaxIdleConns:        numUsers,
			MaxIdleConnsPerHost: numUsers,
		},
		Timeout: 10 * time.Second,
	}

	startBarrier := make(chan struct{})

	utils.InfoLog.Printf("[Client Engine] Initializing %d mock users...", numUsers)
	wg.Add(numUsers)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= numUsers; i++ {
		go func(id int, userRng int) {
			defer wg.Done()

			userID := fmt.Sprintf("user_%04d", id)
			answer := "yes"
			if userRng < 2 {
				// roughly 20% users will supply incorrect response organically
				answer = "no"
			}

			payload := models.AnswerPayload{
				UserID: userID,
				Answer: answer,
			}

			payloadBytes, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(payloadBytes))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")

			// WAIT BARRIER
			<-startBarrier

			// BLAST API
			resp, err := httpClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
		}(i, rng.Intn(10))
	}

	// Give routines a moment to spin up prior
	time.Sleep(2 * time.Second)

	utils.InfoLog.Printf("[Client Engine] Firing %d concurrent requests instantly...\n", numUsers)
	close(startBarrier)

	wg.Wait()
	utils.InfoLog.Println("[Client Engine] All mock users have finished sending responses.")
}
