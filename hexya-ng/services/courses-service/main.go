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

// Server holds the dependencies for the service.
type Server struct {
	coursesRS orm.RecordSetI
}

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

	server := &Server{
		coursesRS: coursesRS,
	}

	http.HandleFunc("/courses", server.coursesHandler)
	http.HandleFunc("/courses/", server.courseHandler)
	fmt.Println("Courses service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func (s *Server) coursesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.listCourses(w, r)
	case "POST":
		s.createCourse(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) courseHandler(w http.ResponseWriter, r *http.Request) {
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
		s.getCourse(w, r, id)
	case "PUT":
		s.updateCourse(w, r, id)
	case "DELETE":
		s.deleteCourse(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := s.coursesRS.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (s *Server) createCourse(w http.ResponseWriter, r *http.Request) {
	var courseData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&courseData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCourse, err := s.coursesRS.Create(courseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCourse)
}

func (s *Server) getCourse(w http.ResponseWriter, r *http.Request, id int) {
	course, err := s.coursesRS.ReadOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func (s *Server) updateCourse(w http.ResponseWriter, r *http.Request, id int) {
	var courseData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&courseData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.coursesRS.Write(id, courseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteCourse(w http.ResponseWriter, r *http.Request, id int) {
	if err := s.coursesRS.Unlink(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
