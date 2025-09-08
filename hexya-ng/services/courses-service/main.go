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
	orm.NewModel("Course")
	orm.GetModel("Course").AddFields(map[string]*orm.FieldInfo{
		"Name":        {Name: "Name", Type: "Char", String: "Name", Required: true},
		"Description": {Name: "Description", Type: "Char", String: "Description"},
	})
}

func main() {
	// Database connection
	orm.Init("sqlite3", "./courses.db")
	defer orm.Close()

	// Create table if it doesn't exist.
	coursesRS := orm.NewRecordSet("Course")
	if err := coursesRS.CreateTable(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

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
	path := strings.TrimPrefix(r.URL.Path, "/courses/")
	if path == "" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(path)
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
	coursesRS := orm.NewRecordSet("Course")
	courses, err := coursesRS.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var courseData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&courseData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coursesRS := orm.NewRecordSet("Course")
	newCourse, err := coursesRS.Create(courseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCourse)
}

func getCourse(w http.ResponseWriter, r *http.Request, id int) {
	// The basic ORM doesn't support reading a single record yet.
	// This is a placeholder.
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func updateCourse(w http.ResponseWriter, r *http.Request, id int) {
	// The basic ORM doesn't support updating yet.
	// This is a placeholder.
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func deleteCourse(w http.ResponseWriter, r *http.Request, id int) {
	// The basic ORM doesn't support deleting yet.
	// This is a placeholder.
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
