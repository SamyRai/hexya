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

// Session represents a session in the system.
type Session struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CourseID  int    `json:"course_id"`
	StartDate string `json:"start_date"`
	Duration  int    `json:"duration"`
	Seats     int    `json:"seats"`
}

var (
	sessions     = make(map[int]Session)
	nextID       = 1
	sessionsMutex = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/sessions", sessionsHandler)
	http.HandleFunc("/sessions/", sessionHandler)
	http.HandleFunc("/courses/", courseSessionsHandler)
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
	id, err := strconv.Atoi(r.URL.Path[len("/sessions/"):])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getSession(w, r, id)
	case "PUT":
		updateSession(w, r, id)
	case "DELETE":
		deleteSession(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func courseSessionsHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	// Expected path: /courses/{course_id}/sessions
	if len(pathParts) != 4 || pathParts[3] != "sessions" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		listSessionsByCourse(w, r, courseID)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listSessions(w http.ResponseWriter, r *http.Request) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessionList := make([]Session, 0, len(sessions))
	for _, session := range sessions {
		sessionList = append(sessionList, session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionList)
}

func createSession(w http.ResponseWriter, r *http.Request) {
	var session Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session.ID = nextID
	nextID++
	sessions[session.ID] = session

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func getSession(w http.ResponseWriter, r *http.Request, id int) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, ok := sessions[id]
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func updateSession(w http.ResponseWriter, r *http.Request, id int) {
	var updatedSession Session
	if err := json.NewDecoder(r.Body).Decode(&updatedSession); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	_, ok := sessions[id]
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	updatedSession.ID = id
	sessions[id] = updatedSession

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSession)
}

func deleteSession(w http.ResponseWriter, r *http.Request, id int) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	_, ok := sessions[id]
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	delete(sessions, id)
	w.WriteHeader(http.StatusNoContent)
}

func listSessionsByCourse(w http.ResponseWriter, r *http.Request, courseID int) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessionList := make([]Session, 0)
	for _, session := range sessions {
		if session.CourseID == courseID {
			sessionList = append(sessionList, session)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionList)
}
