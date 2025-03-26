package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TriggerEventHandler handles "/trigger-event" to demonstrate an event-based action
func TriggerEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Example: run a goroutine for asynchronous work
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("Asynchronous event triggered!")
	}()

	resp := map[string]string{"status": "Event triggered successfully!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
