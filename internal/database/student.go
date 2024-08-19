package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang-assignment/internal/student"

	log "github.com/sirupsen/logrus"
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

func (d *StudentStore) PostStudent(ctx context.Context, stud student.Student) (student.Student, error) {

	log.Printf("Creating student: CreatedBy=%s, UpdatedBy=%s", stud.CreatedBy, stud.UpdatedBy)

	stud.CreatedOn = time.Now()
	stud.UpdatedOn = stud.CreatedOn
	_, err := d.DB.NamedExecContext(ctx, `INSERT INTO students (id, created_by, created_on, updated_by, updated_on, name, email, age, course)
        VALUES (:id, :created_by, :created_on, :updated_by, :updated_on, :name, :email, :age, :course)`,
		stud)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert student: %w", err)
	}
	return stud, nil
}

func (d *StudentStore) UpdateStudent(ctx context.Context, id string, stud student.Student) (student.Student, error) {

	if id != stud.ID {
		return student.Student{}, fmt.Errorf("mismatching student ID")
	}

	query := `UPDATE students SET
		created_by = :created_by,
        updated_by = :updated_by,
        updated_on = :updated_on,
        name = :name,
        email = :email,
        age = :age,
        course = :course
        WHERE id = :id`

	result, err := d.DB.NamedExecContext(ctx, query, stud)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to update student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return student.Student{}, fmt.Errorf("could not determine rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return student.Student{}, fmt.Errorf("no rows were updated, student with ID %s might not exist", id)
	}

	return stud, nil
}

func (s *StudentStore) DeleteStudent(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	return nil
}
