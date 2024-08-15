package transport

import (
	"context"
	"encoding/json"
	"errors"
	"golang-assignment/internal/student"
	"net/http"
	"time"

	util "golang-assignment/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// StudentService defines the interface for student operations
type StudentService interface {
	GetStudent(ctx context.Context, ID string) (student.Student, error)
	PostStudent(ctx context.Context, stu student.Student) (student.Student, error)
	UpdateStudent(ctx context.Context, ID string, newStu student.Student) (student.Student, error)
	DeleteStudent(ctx context.Context, ID string) error
	ReadyCheck(ctx context.Context) error
}

// GetStudent - retrieves a student by ID
func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	stu, err := h.Service.GetStudent(r.Context(), id)
	if err != nil {
		if errors.Is(err, student.ErrFetchingStudent) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(stu); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// PostStudentRequest defines the request body for creating a student
type PostStudentRequest struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Age       int       `json:"age" validate:"required,gt=0"`
	Course    string    `json:"course" validate:"required"`
	CreatedBy string    `json:"created_by"` // Optional or auto-managed
	CreatedOn time.Time `json:"created_on"` // Optional or auto-managed
	UpdatedBy string    `json:"updated_by"` // Optional or auto-managed
	UpdatedOn time.Time `json:"updated_on"` // Optional or auto-managed
}

// studentFromPostStudentRequest converts PostStudentRequest to student.Student
func studentFromPostStudentRequest(u PostStudentRequest) student.Student {
	return student.Student{
		ID:     u.ID, // Generate or handle ID
		Name:   u.Name,
		Email:  u.Email,
		Age:    u.Age,
		Course: u.Course,
	}
}

// PostStudent - adds a new student
func (h *Handler) PostStudent(w http.ResponseWriter, r *http.Request) {
	var postStuReq PostStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&postStuReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(postStuReq); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	stu := studentFromPostStudentRequest(postStuReq)
	stu.CreatedBy = util.GetCurrentUserID(r.Context()) // Assuming util.GetCurrentUserID is correctly defined
	stu.UpdatedBy = stu.CreatedBy
	stu.CreatedOn = time.Now().UTC()
	stu.UpdatedOn = stu.CreatedOn
	stu, err := h.Service.PostStudent(r.Context(), stu)
	if err != nil {
		log.Error(err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(stu); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateStudentRequest defines the request body for updating a student
type UpdateStudentRequest struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Age    int    `json:"age" validate:"required,gt=0"`
	Course string `json:"course" validate:"required"`
}

// studentFromUpdateStudentRequest converts UpdateStudentRequest to student.Student
func studentFromUpdateStudentRequest(u UpdateStudentRequest) student.Student {
	return student.Student{
		Name:   u.Name,
		Email:  u.Email,
		Age:    u.Age,
		Course: u.Course,
	}
}

// UpdateStudent - updates a student by ID
func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	var updateStuRequest UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStuRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(updateStuRequest); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	stu := studentFromUpdateStudentRequest(updateStuRequest)

	stu, err := h.Service.UpdateStudent(r.Context(), studentID, stu)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(stu); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteStudent - deletes a student by ID
func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Successfully Deleted"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
