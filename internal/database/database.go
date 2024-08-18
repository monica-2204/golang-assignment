package database

import (
	"context"
	"fmt"
	"golang-assignment/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// InitDatabase initializes the database connection using the provided configuration.
func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	var err error
	// Initialize the db variable with a connection
	db, err = sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// Ping checks if the database is reachable.
func (store *StudentStore) Ping(ctx context.Context) error {
	return store.DB.PingContext(ctx)
}

// StudentStore - implements the student.StudentStore interface
type StudentStore struct {
	DB *sqlx.DB
}

// NewStudentStore - returns a new instance of StudentStore
func NewStudentStore(db *sqlx.DB) *StudentStore {
	return &StudentStore{DB: db}
}
