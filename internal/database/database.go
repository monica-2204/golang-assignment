package database

import (
	"context"
	"fmt"
	"golang-assignment/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	var err error

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

func (store *StudentStore) Ping(ctx context.Context) error {
	return store.DB.PingContext(ctx)
}

type StudentStore struct {
	DB *sqlx.DB
}

func NewStudentStore(db *sqlx.DB) *StudentStore {
	return &StudentStore{DB: db}
}
