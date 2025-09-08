package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hexya-ng/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRecordSet is a mock implementation of the RecordSetI interface.
type MockRecordSet struct {
	orm.RecordSet
	CreateFunc  func(data map[string]interface{}) (map[string]interface{}, error)
	ReadFunc    func() ([]map[string]interface{}, error)
	ReadOneFunc func(id int) (map[string]interface{}, error)
	WriteFunc   func(id int, data map[string]interface{}) error
	UnlinkFunc  func(id int) error
}

func (m *MockRecordSet) Create(data map[string]interface{}) (map[string]interface{}, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(data)
	}
	return nil, nil
}

func (m *MockRecordSet) Read() ([]map[string]interface{}, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc()
	}
	return nil, nil
}

func (m *MockRecordSet) ReadOne(id int) (map[string]interface{}, error) {
	if m.ReadOneFunc != nil {
		return m.ReadOneFunc(id)
	}
	return nil, nil
}

func (m *MockRecordSet) Write(id int, data map[string]interface{}) error {
	if m.WriteFunc != nil {
		return m.WriteFunc(id, data)
	}
	return nil
}

func (m *MockRecordSet) Unlink(id int) error {
	if m.UnlinkFunc != nil {
		return m.UnlinkFunc(id)
	}
	return nil
}

func (m *MockRecordSet) CreateTable() error {
	return nil
}

func TestListCourses(t *testing.T) {
	mockRS := &MockRecordSet{
		ReadFunc: func() ([]map[string]interface{}, error) {
			return []map[string]interface{}{
				{"Name": "Course 1"},
				{"Name": "Course 2"},
			}, nil
		},
	}

	server := &Server{coursesRS: mockRS}

	req, err := http.NewRequest("GET", "/courses", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.coursesHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var courses []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &courses)
	require.NoError(t, err)

	assert.Len(t, courses, 2)
	assert.Equal(t, "Course 1", courses[0]["Name"])
}

func TestCreateCourse(t *testing.T) {
	courseData := map[string]interface{}{"Name": "New Course"}
	mockRS := &MockRecordSet{
		CreateFunc: func(data map[string]interface{}) (map[string]interface{}, error) {
			return courseData, nil
		},
	}

	server := &Server{coursesRS: mockRS}

	body, err := json.Marshal(courseData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/courses", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.coursesHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var newCourse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &newCourse)
	require.NoError(t, err)

	assert.Equal(t, "New Course", newCourse["Name"])
}

func TestGetCourse(t *testing.T) {
	mockRS := &MockRecordSet{
		ReadOneFunc: func(id int) (map[string]interface{}, error) {
			if id == 1 {
				return map[string]interface{}{"Name": "Course 1"}, nil
			}
			return nil, nil
		},
	}

	server := &Server{coursesRS: mockRS}

	req, err := http.NewRequest("GET", "/courses/1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.courseHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var course map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &course)
	require.NoError(t, err)

	assert.Equal(t, "Course 1", course["Name"])
}

func TestUpdateCourse(t *testing.T) {
	updateData := map[string]interface{}{"Name": "Updated Course"}
	mockRS := &MockRecordSet{
		WriteFunc: func(id int, data map[string]interface{}) error {
			assert.Equal(t, 1, id)
			assert.Equal(t, updateData, data)
			return nil
		},
	}

	server := &Server{coursesRS: mockRS}

	body, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/courses/1", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.courseHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteCourse(t *testing.T) {
	mockRS := &MockRecordSet{
		UnlinkFunc: func(id int) error {
			assert.Equal(t, 1, id)
			return nil
		},
	}

	server := &Server{coursesRS: mockRS}

	req, err := http.NewRequest("DELETE", "/courses/1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.courseHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
