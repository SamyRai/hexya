package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hexya-ng/orm"
)

func init() {
	orm.NewModel("Attendee")
	orm.GetModel("Attendee").AddFields(map[string]*orm.FieldInfo{
		"Name":      {Name: "Name", Type: "Char", String: "Name", Required: true},
		"SessionID": {Name: "SessionID", Type: "Integer", String: "Session ID"},
	})
}

func main() {
	// Database connection
	orm.Init("sqlite3", "./attendees.db")
	defer orm.Close()

	// Create table if it doesn't exist.
	attendeesRS := orm.NewRecordSet("Attendee")
	if err := attendeesRS.CreateTable(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/attendees", attendeesHandler)
	http.HandleFunc("/attendees/", attendeeHandler)
	// I will handle the /sessions/{id}/attendees endpoint later.
	// http.HandleFunc("/sessions/", sessionAttendeesHandler)
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
	path := strings.TrimPrefix(r.URL.Path, "/attendees/")
	if path == "" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid attendee ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getAttendee(w, r, id)
	// Other methods are not implemented yet in the ORM.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listAttendees(w http.ResponseWriter, r *http.Request) {
	attendeesRS := orm.NewRecordSet("Attendee")
	attendees, err := attendeesRS.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendees)
}

func createAttendee(w http.ResponseWriter, r *http.Request) {
	var attendeeData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&attendeeData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if session exists
	sessionID, ok := attendeeData["SessionID"].(float64) // JSON numbers are float64
	if !ok {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionURL := fmt.Sprintf("http://localhost:8082/sessions/%d", int(sessionID))
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

	attendeesRS := orm.NewRecordSet("Attendee")
	newAttendee, err := attendeesRS.Create(attendeeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAttendee)
}

func getAttendee(w http.ResponseWriter, r *http.Request, id int) {
	// The basic ORM doesn't support reading a single record yet.
	// This is a placeholder.
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
