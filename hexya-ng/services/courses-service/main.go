package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Course represents a course in the system.
type Course struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

var db *gorm.DB

func main() {
	// Database connection
	var err error
	db, err = gorm.Open(sqlite.Open("courses.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&Course{})

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
	var courses []Course
	if err := db.Find(&courses).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Create(&course).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

func getCourse(w http.ResponseWriter, r *http.Request, id int) {
	var course Course
	if err := db.First(&course, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func updateCourse(w http.ResponseWriter, r *http.Request, id int) {
	var course Course
	if err := db.First(&course, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updatedCourse Course
	if err := json.NewDecoder(r.Body).Decode(&updatedCourse); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	course.Name = updatedCourse.Name
	course.Description = updatedCourse.Description

	if err := db.Save(&course).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func deleteCourse(w http.ResponseWriter, r *http.Request, id int) {
	var course Course
	if err := db.First(&course, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := db.Delete(&course).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
