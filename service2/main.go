package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Request struct to receive user input
type Request struct {
	Date string `json:"date"`
}

// Response struct to send messages to frontend
type Response struct {
	Message string `json:"message"`
}

// Store active date tracking goroutines and WebSocket connections
var (
	requestedDates = make(map[string]bool)
	clients        = make(map[*websocket.Conn]bool)
	mutex          sync.Mutex
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Handle incoming WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Keep connection open
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
	}
}

// Handle user requests to set a date
func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	if _, exists := requestedDates[req.Date]; !exists {
		requestedDates[req.Date] = true
		go waitForDate(req.Date) // Start independent tracking
	}
	mutex.Unlock()

	response := Response{Message: fmt.Sprintf("Date %s has been set!", req.Date)}
	json.NewEncoder(w).Encode(response)
}

// Goroutine to wait for the requested date independently
func waitForDate(requestedDate string) {
	for {
		currentDate := time.Now().Format("2006-01-02")

		if currentDate == requestedDate {
			message := fmt.Sprintf("ALERT: Today (%s) is the requested day!", requestedDate)
			fmt.Println(message) // Print to server console
			sendToAllClients(message) // Send to all WebSocket clients

			mutex.Lock()
			delete(requestedDates, requestedDate)
			mutex.Unlock()
			return
		}

		time.Sleep(1 * time.Hour) // Check every hour
	}
}

// Send message to all WebSocket clients
func sendToAllClients(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(Response{Message: message})
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {
	// Serve frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// API endpoint to set dates
	http.HandleFunc("/check-date", handleRequest)

	// WebSocket connection for real-time updates
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server started at 8083")
	http.ListenAndServe(":8083", nil)
}
