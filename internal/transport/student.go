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
	AuthenticateUser(ctx context.Context, userID, password string) (student.User, error) // Update here
	GenerateJWT(user student.User) (string, error)
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
	CreatedBy string    `json:"created_by"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedOn time.Time `json:"updated_on"`
}

// studentFromPostStudentRequest converts PostStudentRequest to student.Student
func studentFromPostStudentRequest(u PostStudentRequest) student.Student {
	return student.Student{
		ID:     u.ID,
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
	stu.CreatedBy = util.GetCurrentUserID(r.Context())
	stu.UpdatedBy = stu.CreatedBy
	log.Printf("userID: %s", stu.CreatedBy)

	stu.CreatedOn = time.Now()
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

	// Validate the incoming request
	validate := validator.New()
	if err := validate.Struct(updateStuRequest); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	existingStudent, err := h.Service.GetStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	// Convert the request to a Student struct
	stu := studentFromUpdateStudentRequest(updateStuRequest)
	stu.CreatedBy = existingStudent.CreatedBy
	stu.ID = studentID // Ensure the ID is set for updating the correct student
	stu.UpdatedBy = util.GetCurrentUserID(r.Context())
	stu.UpdatedOn = time.Now()

	// Attempt to update the student
	updatedStu, err := h.Service.UpdateStudent(r.Context(), studentID, stu)
	if err != nil {
		log.Printf("Error updating student: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	// Respond with the updated student information
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedStu); err != nil {
		log.Printf("Error encoding response: %v", err)
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

	// Check if the student exists
	_, err := h.Service.GetStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Proceed to delete the student
	err = h.Service.DeleteStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Successfully Deleted"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
