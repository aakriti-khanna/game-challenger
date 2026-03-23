package main

import (
	"net/http"
	"time"

	"game-challenge/client"
	"game-challenge/engine"
	"game-challenge/routes"
	"game-challenge/utils"
)

func main() {
	// 0. Start our asynchronous Event-Listener channel receiver
	engine.StartEngine()

	// 1. Setup global multiplexer and bind endpoints
	mux := http.NewServeMux()
	routes.SetupRoutes(mux)

	// 2. Wrap via native Go Server struct
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 3. Initiate the API thread asynchronously
	go func() {
		utils.InfoLog.Println("[API Server] Starting on port :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.ErrorLog.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Temporary pause to guarantee socket binds locally before mocking traffic
	time.Sleep(1 * time.Second)

	// 4. Trigger High Velocity Mock Engine on API Server Instance
	client.MockUserEngine(5000, "http://localhost:8080/api/submit")

	// 5. Signal the Event-Channel to stop taking requests. Wait for pending queue to print metrics gracefully
	engine.StopEngine()

	utils.InfoLog.Println("Simulation completely done. Exiting.")
}
