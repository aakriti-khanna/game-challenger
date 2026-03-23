package engine

import (
	"sync"
	"sync/atomic"
	"time"

	"game-challenge/utils"
)

var (
	winnerFound int32
	winnerID    string
	winningTime time.Time

	winnerOnce sync.Once

	// Metrics trackers
	totalRequests  uint64
	correctAnswers uint64
	wrongAnswers   uint64
)

// RecordAnswer safely increments our metrics in a concurrent manner identically to checking the winner lock
func RecordAnswer(answer string) {
	atomic.AddUint64(&totalRequests, 1) // Increment total payloads parsed
	if answer == "yes" {
		atomic.AddUint64(&correctAnswers, 1) // Increment correct answers
	} else {
		atomic.AddUint64(&wrongAnswers, 1) // Increment wrong answers
	}
}

// CheckAndRecordWinner atomically locks in the first person who gets the correct answer.
func CheckAndRecordWinner(userID string, receiveTime time.Time) {
	// Only attempt lock if game continues
	if atomic.LoadInt32(&winnerFound) == 1 {
		return
	}

	winnerOnce.Do(func() {
		atomic.StoreInt32(&winnerFound, 1) // Mark winner as found safely
		winnerID = userID
		winningTime = receiveTime

		// Console output
		utils.LogWinner(winnerID, winningTime.Format("2006-01-02 15:04:05.000000"))
	})
}

// IsGameOver returns if a winner has already been found
func IsGameOver() bool {
	return atomic.LoadInt32(&winnerFound) == 1
}

// PrintMetrics logs out the aggregate results
func PrintMetrics() {
	utils.LogMetrics(
		atomic.LoadUint64(&totalRequests),
		atomic.LoadUint64(&correctAnswers),
		atomic.LoadUint64(&wrongAnswers),
	)
}
