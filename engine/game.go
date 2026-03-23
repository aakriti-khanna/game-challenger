package engine

import (
	"time"

	"game-challenge/utils"
)

type AnswerEvent struct {
	UserID      string
	Answer      string
	ReceiveTime time.Time
}

var (
	eventsChan = make(chan AnswerEvent, 10000)

	doneChan = make(chan struct{})

	winnerFound bool
	winnerID    string
	winningTime time.Time

	totalRequests  uint64
	correctAnswers uint64
	wrongAnswers   uint64
)

func StartEngine() {
	go func() {

		for event := range eventsChan {
			totalRequests++
			if event.Answer == "yes" {
				correctAnswers++

				if !winnerFound {
					winnerFound = true
					winnerID = event.UserID
					winningTime = event.ReceiveTime

					utils.LogWinner(winnerID, winningTime.Format("2006-01-02 15:04:05.000000"))
				}
			} else {
				wrongAnswers++
			}
		}
		// Unblock the closer
		close(doneChan)
	}()
}

// ProcessEvent pushes HTTP traffic linearly onto the engine queue instantly
func ProcessEvent(userID, answer string, receiveTime time.Time) {
	eventsChan <- AnswerEvent{
		UserID:      userID,
		Answer:      answer,
		ReceiveTime: receiveTime,
	}
}

// StopEngine gracefully terminates the channel and prints the recorded stats
func StopEngine() {
	// Signal no more events will natively arrive
	close(eventsChan)

	// Blocks natively until the event loop function (goroutine) completes its channel queue completely
	<-doneChan

	utils.LogMetrics(totalRequests, correctAnswers, wrongAnswers)
}
