package database

import (
	"context"
	"database/sql"
	"errors"
	"golang-assignment/internal/student"
)

// StudentStore - implements the student.StudentStore interface
type StudentStore struct {
	DB *sql.DB
}

// NewStudentStore - returns a new instance of StudentStore
func NewStudentStore(db *sql.DB) *StudentStore {
	return &StudentStore{DB: db}
}

// GetStudent - retrieves a student by ID from the database
func (store *StudentStore) GetStudent(ctx context.Context, id string) (student.Student, error) {
	var s student.Student
	query := `SELECT id, created_by, created_on, updated_by, updated_on, name, email, age, course FROM students WHERE id = ?`

	row := store.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&s.ID, &s.CreatedBy, &s.CreatedOn, &s.UpdatedBy, &s.UpdatedOn, &s.Name, &s.Email, &s.Age, &s.Course)
	if err != nil {
		if err == sql.ErrNoRows {
			return s, errors.New("student with id " + id + " not found")
		}
		return s, errors.New("error fetching student with id " + id + ": " + err.Error())
	}

	return s, nil
}

// PostStudent - adds a new student to the database
func (store *StudentStore) PostStudent(ctx context.Context, s student.Student) (student.Student, error) {
	query := `INSERT INTO students (id, created_by, created_on, updated_by, updated_on, name, email, age, course) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := store.DB.ExecContext(ctx, query, s.ID, s.CreatedBy, s.CreatedOn, s.UpdatedBy, s.UpdatedOn, s.Name, s.Email, s.Age, s.Course)
	if err != nil {
		return s, errors.New("error inserting student: " + err.Error())
	}

	return s, nil
}

// UpdateStudent - updates an existing student in the database
func (store *StudentStore) UpdateStudent(ctx context.Context, id string, s student.Student) (student.Student, error) {
	query := `UPDATE students SET name = ?, email = ?, age = ?, course = ?, updated_by = ?, updated_on = ? WHERE id = ?`

	_, err := store.DB.ExecContext(ctx, query, s.Name, s.Email, s.Age, s.Course, s.UpdatedBy, s.UpdatedOn, id)
	if err != nil {
		return s, errors.New("error updating student with id " + id + ": " + err.Error())
	}

	return s, nil
}

// DeleteStudent - deletes a student from the database by ID
func (store *StudentStore) DeleteStudent(ctx context.Context, id string) error {
	query := `DELETE FROM students WHERE id = ?`

	result, err := store.DB.ExecContext(ctx, query, id)
	if err != nil {
		return errors.New("error deleting student with id " + id + ": " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("error checking rows affected when deleting student with id " + id + ": " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("no student found with id " + id + " to delete")
	}

	return nil
}

// Ping - checks if the database is reachable
func (store *StudentStore) Ping(ctx context.Context) error {
	if err := store.DB.PingContext(ctx); err != nil {
		return errors.New("database ping failed: " + err.Error())
	}
	return nil
}
