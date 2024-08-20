package main

import (
	"context"
	"golang-assignment/config"
	"golang-assignment/internal/database"
	"golang-assignment/internal/student"
	"golang-assignment/internal/transport"
	"os"

	log "github.com/sirupsen/logrus"
)

func Run() error {

	log.Info("Setting up our application")

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.JSONFormatter{})

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load configuration")
		return err
	}

	// Initialize database connection
	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Error("failed to setup connection to the database")
		return err
	}

	// Initialize the student store and service
	studentStore := database.NewStudentStore(db)
	studentService := student.NewService(studentStore)

	ctx := context.Background()
	if err := studentStore.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	// Initialize the HTTP handler
	handler := transport.NewHandler(studentService)

	// Start the HTTP server
	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")

	}
}
