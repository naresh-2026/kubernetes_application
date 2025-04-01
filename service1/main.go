package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"reflect"

	"github.com/gorilla/websocket"
)

var scheduledTime string
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var clients = make(map[*websocket.Conn]bool) // Store WebSocket clients

// Serve the HTML page from static directory
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

// WebSocket connection handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	fmt.Println("New WebSocket Client Connected!")

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket Client Disconnected")
			delete(clients, conn)
			break
		}
	}
}

// Handle form submission with JSON request body
func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Time string `json:"time"`
	}

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	scheduledTime = data.Time
	go task(scheduledTime)

	// Send response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Time scheduled successfully"})
}

// Background task to check and notify when the time matches
func task(setTime string) {
	fmt.Println("Processing request....",setTime);
	for {
		fmt.Println("checking")
		fmt.Println(reflect.TypeOf(setTime))
		fmt.Println(reflect.TypeOf(time.Now().Format("15:04")))
		fmt.Println(setTime)
		fmt.Println(time.Now().Format("15:04"))
		loc, _ := time.LoadLocation("Local") // Load the system's local timezone
		adjustedTime := time.Now().In(loc).Add(5*time.Hour + 30*time.Minute).Format("15:04")
		fmt.Println("Adjusted Time:", adjustedTime)

		if adjustedTime == setTime {
        //if time.Now().In(loc).Format("15:04") == setTime {
			fmt.Println("inside")
			message := fmt.Sprintf("Notify the user: It's %s", setTime)
			fmt.Println(message)
			notifyClients(message)
			break
		}
		time.Sleep(30 * time.Second) // Check every 30 seconds
	}
}

// Send WebSocket message to all connected clients
func notifyClients(message string) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("WebSocket Write Error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {
	fmt.Println("Server started on port 8082...")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/schedule", scheduleHandler)
	http.HandleFunc("/ws", wsHandler) // WebSocket endpoint

	http.ListenAndServe(":8082", nil)
}
