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

	// Capture response time immediately for checking who was strictly first later
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

	// 1. Record metrics: Even if the game is "won", we must record the payload info
	engine.RecordAnswer(payload.Answer)

	// 2. Only allow correct answers to attempt obtaining the winner lock
	if payload.Answer == "yes" {
		if !engine.IsGameOver() {
			engine.CheckAndRecordWinner(payload.UserID, receiveTime)
		}
	}

	w.WriteHeader(http.StatusOK)
	if engine.IsGameOver() {
		w.Write([]byte("Answer received. Game is already over."))
	} else {
		w.Write([]byte("Answer received."))
	}
}
