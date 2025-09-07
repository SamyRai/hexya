package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Course represents a course in the system.
type Course struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	courses  = make(map[int]Course)
	nextID   = 1
	coursesMutex = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/courses", coursesHandler)
	http.HandleFunc("/courses/", courseHandler)
	fmt.Println("Courses service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listCourses(w, r)
	case "POST":
		createCourse(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/courses/"):])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getCourse(w, r, id)
	case "PUT":
		updateCourse(w, r, id)
	case "DELETE":
		deleteCourse(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listCourses(w http.ResponseWriter, r *http.Request) {
	coursesMutex.Lock()
	defer coursesMutex.Unlock()

	courseList := make([]Course, 0, len(courses))
	for _, course := range courses {
		courseList = append(courseList, course)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courseList)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coursesMutex.Lock()
	defer coursesMutex.Unlock()

	course.ID = nextID
	nextID++
	courses[course.ID] = course

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

func getCourse(w http.ResponseWriter, r *http.Request, id int) {
	coursesMutex.Lock()
	defer coursesMutex.Unlock()

	course, ok := courses[id]
	if !ok {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func updateCourse(w http.ResponseWriter, r *http.Request, id int) {
	var updatedCourse Course
	if err := json.NewDecoder(r.Body).Decode(&updatedCourse); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coursesMutex.Lock()
	defer coursesMutex.Unlock()

	_, ok := courses[id]
	if !ok {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	updatedCourse.ID = id
	courses[id] = updatedCourse

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCourse)
}

func deleteCourse(w http.ResponseWriter, r *http.Request, id int) {
	coursesMutex.Lock()
	defer coursesMutex.Unlock()

	_, ok := courses[id]
	if !ok {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	delete(courses, id)
	w.WriteHeader(http.StatusNoContent)
}
