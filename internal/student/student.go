package student

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrFetchingStudent   = errors.New("could not fetch student by ID")
	ErrUpdatingStudent   = errors.New("could not update student")
	ErrNoStudentFound    = errors.New("no student found")
	ErrDeletingStudent   = errors.New("could not delete student")
	ErrListingStudents   = errors.New("could not list students")
	ErrSearchingStudents = errors.New("could not search students")
	ErrNotImplemented    = errors.New("not implemented")
)

type Student struct {
	ID        string    `json:"id" db:"id"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedOn time.Time `json:"created_on" db:"created_on"`
	UpdatedBy string    `json:"updated_by" db:"updated_by"`
	UpdatedOn time.Time `json:"updated_on" db:"updated_on"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Age       int       `json:"age" db:"age"`
	Course    string    `json:"course" db:"course"`
}

type StudentStore interface {
	GetStudent(context.Context, string) (Student, error)
	PostStudent(context.Context, Student) (Student, error)
	UpdateStudent(context.Context, string, Student) (Student, error)
	DeleteStudent(context.Context, string) error
	Ping(context.Context) error
}

type Service struct {
	Store StudentStore
}

func NewService(store StudentStore) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetStudent(ctx context.Context, ID string) (Student, error) {
	student, err := s.Store.GetStudent(ctx, ID)
	if err != nil {
		log.Errorf("an error occurred fetching the student: %s", err.Error())
		return Student{}, ErrFetchingStudent
	}
	return student, nil
}

func (s *Service) PostStudent(ctx context.Context, student Student) (Student, error) {
	student, err := s.Store.PostStudent(ctx, student)
	if err != nil {
		log.Errorf("an error occurred adding the student: %s", err.Error())
	}
	return student, nil
}

func (s *Service) UpdateStudent(
	ctx context.Context, ID string, newStudent Student,
) (Student, error) {
	student, err := s.Store.UpdateStudent(ctx, ID, newStudent)
	if err != nil {
		log.Errorf("an error occurred updating the student: %s", err.Error())
	}
	return student, nil
}

func (s *Service) DeleteStudent(ctx context.Context, ID string) error {
	err := s.Store.DeleteStudent(ctx, ID)
	if err != nil {
		log.Errorf("an error occurred deleting the student: %s", err.Error())
	}
	return err
}

func (s *Service) ReadyCheck(ctx context.Context) error {
	log.Info("Checking readiness")
	return s.Store.Ping(ctx)
}
