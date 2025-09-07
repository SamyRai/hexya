package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Attendee represents an attendee in the system.
type Attendee struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SessionID int    `json:"session_id"`
}

var (
	attendees     = make(map[int]Attendee)
	nextID        = 1
	attendeesMutex = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/attendees", attendeesHandler)
	http.HandleFunc("/attendees/", attendeeHandler)
	http.HandleFunc("/sessions/", sessionAttendeesHandler)
	fmt.Println("Attendees service listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func attendeesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listAttendees(w, r)
	case "POST":
		createAttendee(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func attendeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/attendees/"):])
	if err != nil {
		http.Error(w, "Invalid attendee ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getAttendee(w, r, id)
	case "PUT":
		updateAttendee(w, r, id)
	case "DELETE":
		deleteAttendee(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func sessionAttendeesHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	// Expected path: /sessions/{session_id}/attendees
	if len(pathParts) != 4 || pathParts[3] != "attendees" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	sessionID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		listAttendeesBySession(w, r, sessionID)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listAttendees(w http.ResponseWriter, r *http.Request) {
	attendeesMutex.Lock()
	defer attendeesMutex.Unlock()

	attendeeList := make([]Attendee, 0, len(attendees))
	for _, attendee := range attendees {
		attendeeList = append(attendeeList, attendee)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendeeList)
}

func createAttendee(w http.ResponseWriter, r *http.Request) {
	var attendee Attendee
	if err := json.NewDecoder(r.Body).Decode(&attendee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if session exists
	sessionURL := fmt.Sprintf("http://localhost:8082/sessions/%d", attendee.SessionID)
	resp, err := http.Get(sessionURL)
	if err != nil {
		http.Error(w, "Error verifying session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	attendeesMutex.Lock()
	defer attendeesMutex.Unlock()

	attendee.ID = nextID
	nextID++
	attendees[attendee.ID] = attendee

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attendee)
}

func getAttendee(w http.ResponseWriter, r *http.Request, id int) {
	attendeesMutex.Lock()
	defer attendeesMutex.Unlock()

	attendee, ok := attendees[id]
	if !ok {
		http.Error(w, "Attendee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendee)
}

func updateAttendee(w http.ResponseWriter, r *http.Request, id int) {
	var updatedAttendee Attendee
	if err := json.NewDecoder(r.Body).Decode(&updatedAttendee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	attendeesMutex.Lock()
	defer attendeesMutex.Unlock()

	_, ok := attendees[id]
	if !ok {
		http.Error(w, "Attendee not found", http.StatusNotFound)
		return
	}

	updatedAttendee.ID = id
	attendees[id] = updatedAttendee

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAttendee)
}

func deleteAttendee(w http.ResponseWriter, r *http.Request, id int) {
	attendeesMutex.Lock()
	defer attendeesMutex.Unlock()

	_, ok := attendees[id]
	if !ok {
		http.Error(w, "Attendee not found", http.StatusNotFound)
		return
	}

	delete(attendees, id)
	w.WriteHeader(http.StatusNoContent)
}

func listAttendeesBySession(w http.ResponseWriter, r *http.Request, sessionID int) {
	attendeesMutex.Lock()
	defer attendeesMutex.Unlock()

	attendeeList := make([]Attendee, 0)
	for _, attendee := range attendees {
		if attendee.SessionID == sessionID {
			attendeeList = append(attendeeList, attendee)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendeeList)
}
