package database

import (
	"context"
	"errors"
	"golang-assignment/config"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var db *sqlx.DB

// InitDatabase initializes the database connection using the provided configuration.
func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Connect("mysql", cfg.DatabaseDSN)
	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
		return nil, errors.New("failed to connect to database: " + err.Error())
	}
	return db, nil
}

// Ping checks if the database is reachable.
func Ping(ctx context.Context) error {
	if err := db.PingContext(ctx); err != nil {
		logrus.Errorf("database ping failed: %v", err)
		return errors.New("database ping failed: " + err.Error())
	}
	return nil
}
