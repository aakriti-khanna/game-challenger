package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"game-challenge/engine"
	"game-challenge/models"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/submit", submitAnswerHandler)
}

func submitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Capture response time immediately natively to chronologically lock
	receiveTime := time.Now()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var payload models.AnswerPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Distribute payload asynchronously via Channels directly eliminating all waiting and atomics
	engine.ProcessEvent(payload.UserID, payload.Answer, receiveTime)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event enqueued."))
}
