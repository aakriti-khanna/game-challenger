package utils

import (
	"log"
	"os"
)

var (
	InfoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	GameLog  = log.New(os.Stdout, "🎮 GAME: ", log.Ldate|log.Ltime)
)

func LogWinner(userID string, winningTime string) {
	GameLog.Printf("\n==========================================================================")
	GameLog.Printf("🏆 WINNER DETECTED! First correct answer joined:")
	GameLog.Printf("User ID     : %s", userID)
	GameLog.Printf("Answer Time : %s (microsecond precision)", winningTime)
	GameLog.Printf("==========================================================================\n")
}

// LogMetrics explicitly displays the final statistics collected during the run
func LogMetrics(total uint64, correct uint64, wrong uint64) {
	InfoLog.Printf("\n====================== API METRICS ======================")
	InfoLog.Printf("Total Requests Received   : %d", total)
	InfoLog.Printf("Correct Answers Received  : %d", correct)
	InfoLog.Printf("Incorrect Answers Received: %d", wrong)
	InfoLog.Printf("=========================================================\n")
}
