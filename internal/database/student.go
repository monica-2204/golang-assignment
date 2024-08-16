package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	//"github.com/jmoiron/sqlx"
	"golang-assignment/internal/student"
)

type StudentRow struct {
	ID        string         `db:"id"`
	CreatedBy sql.NullString `db:"created_by"`
	CreatedOn sql.NullTime   `db:"created_on"`
	UpdatedBy sql.NullString `db:"updated_by"`
	UpdatedOn sql.NullTime   `db:"updated_on"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	Age       int            `db:"age"`
	Course    string         `db:"course"`
}

func convertStudentRowToStudent(r StudentRow) student.Student {
	return student.Student{
		ID:        r.ID,
		CreatedBy: r.CreatedBy.String,
		CreatedOn: r.CreatedOn.Time,
		UpdatedBy: r.UpdatedBy.String,
		UpdatedOn: r.UpdatedOn.Time,
		Name:      r.Name,
		Email:     r.Email,
		Age:       r.Age,
		Course:    r.Course,
	}
}

// GetStudent - retrieves a student from the database by ID
func (s *StudentStore) GetStudent(ctx context.Context, id string) (student.Student, error) {
	var studentRow StudentRow
	query := "SELECT id, created_by, created_on, updated_by, updated_on, name, email, age, course FROM students WHERE id = ?"
	err := s.DB.GetContext(ctx, &studentRow, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return student.Student{}, fmt.Errorf("student with ID %s not found", id)
		}
		return student.Student{}, fmt.Errorf("an error occurred fetching the student: %w", err)
	}

	return convertStudentRowToStudent(studentRow), nil
}

// PostStudent - adds a new student to the database
func (d *StudentStore) PostStudent(ctx context.Context, stud student.Student) (student.Student, error) {
	// Insert implementation
	//stud.ID = generateUniqueID()

	stud.CreatedBy = "admin" // Assumes you have a function to get the user ID from the context
	stud.CreatedOn = time.Now().UTC()
	stud.UpdatedBy = stud.CreatedBy
	stud.UpdatedOn = stud.CreatedOn
	_, err := d.DB.NamedExecContext(ctx, `INSERT INTO students (id, created_by, created_on, updated_by, updated_on, name, email, age, course)
        VALUES (:id, :created_by, :created_on, :updated_by, :updated_on, :name, :email, :age, :course)`,
		stud)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert student: %w", err)
	}
	return stud, nil
}

// UpdateStudent - updates an existing student in the database
func (d *StudentStore) UpdateStudent(ctx context.Context, id string, stud student.Student) (student.Student, error) {
	// Ensure the ID matches between the request and the struct
	if id != stud.ID {
		return student.Student{}, fmt.Errorf("mismatching student ID")
	}

	// Updating fields typically meant for creation should not be modified here
	stud.UpdatedBy = "admin" // You might want to fetch this from context or another reliable source
	stud.UpdatedOn = time.Now().UTC()

	// Update only necessary fields
	query := `UPDATE students SET
		created_by = :created_by,
        updated_by = :updated_by,
        updated_on = :updated_on,
        name = :name,
        email = :email,
        age = :age,
        course = :course
        WHERE id = :id`

	// Execute the query
	result, err := d.DB.NamedExecContext(ctx, query, stud)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to update student: %w", err)
	}

	// Check if the update affected any rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return student.Student{}, fmt.Errorf("could not determine rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return student.Student{}, fmt.Errorf("no rows were updated, student with ID %s might not exist", id)
	}

	return stud, nil
}

// DeleteStudent - deletes a student from the database
func (s *StudentStore) DeleteStudent(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	return nil
}
