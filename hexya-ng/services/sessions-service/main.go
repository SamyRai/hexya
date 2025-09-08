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
	orm.NewModel("Session")
	orm.GetModel("Session").AddFields(map[string]*orm.FieldInfo{
		"Name":      {Name: "Name", Type: "Char", String: "Name", Required: true},
		"CourseID":  {Name: "CourseID", Type: "Integer", String: "Course ID"},
		"StartDate": {Name: "StartDate", Type: "Char", String: "Start Date"},
		"Duration":  {Name: "Duration", Type: "Integer", String: "Duration"},
		"Seats":     {Name: "Seats", Type: "Integer", String: "Seats"},
	})
}

func main() {
	// Database connection
	orm.Init("sqlite3", "./sessions.db")
	defer orm.Close()

	// Create table if it doesn't exist.
	sessionsRS := orm.NewRecordSet("Session")
	if err := sessionsRS.CreateTable(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/sessions", sessionsHandler)
	http.HandleFunc("/sessions/", sessionHandler)
	// I will handle the /courses/{id}/sessions endpoint later.
	// http.HandleFunc("/courses/", courseSessionsHandler)
	fmt.Println("Sessions service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func sessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listSessions(w, r)
	case "POST":
		createSession(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sessions/")
	if path == "" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getSession(w, r, id)
	// Other methods are not implemented yet in the ORM.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listSessions(w http.ResponseWriter, r *http.Request) {
	sessionsRS := orm.NewRecordSet("Session")
	sessions, err := sessionsRS.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func createSession(w http.ResponseWriter, r *http.Request) {
	var sessionData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionsRS := orm.NewRecordSet("Session")
	newSession, err := sessionsRS.Create(sessionData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSession)
}

func getSession(w http.ResponseWriter, r *http.Request, id int) {
	// The basic ORM doesn't support reading a single record yet.
	// This is a placeholder.
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
