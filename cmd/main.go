package main

import (
	"golang-assignment/config"
	"golang-assignment/internal/database"
	"golang-assignment/internal/student"
	transportHTTP "golang-assignment/internal/transport"

	logrus "github.com/sirupsen/logrus"
)

// Run - sets up our application
func Run() error {
	// Set up logging format to JSON
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Setting up our application")

	// Load the configuration (assuming this function exists and returns the necessary configuration)
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Error("failed to load configuration")
		return err
	}

	// Initialize database connection
	db, err := database.InitDatabase(cfg)
	if err != nil {
		logrus.Error("failed to setup connection to the database")
		return err
	}

	// Initialize the student store and service
	studentStore := database.NewStudentStore(db)
	studentService := student.NewService(studentStore)

	// Initialize the HTTP handler
	handler := transportHTTP.NewStudentHandler(studentService)

	// Start the HTTP server
	if err := transportHTTP.Serve(handler); err != nil {
		logrus.Error("failed to gracefully serve our application")
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		logrus.Error(err)
		logrus.Fatal("Error starting up our REST API")
	}
}
