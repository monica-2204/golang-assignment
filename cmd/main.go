package main

import (
	"context"
	"fmt"
	"golang-assignment/config"
	"golang-assignment/internal/database"
	"golang-assignment/internal/student"
	"golang-assignment/internal/transport"
	"golang-assignment/utils"

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
	token, err := utils.GenerateJWT("user123") // Replace "user123" with the user ID or relevant data
	if err != nil {
		logrus.Fatalf("Failed to generate JWT: %v", err)
	}
	fmt.Println("Generated JWT Token:", token)
	// Initialize database connection
	db, err := database.InitDatabase(cfg)
	if err != nil {
		logrus.Error("failed to setup connection to the database")
		return err
	}

	// Initialize the student store and service
	studentStore := database.NewStudentStore(db)
	studentService := student.NewService(studentStore)

	ctx := context.Background()
	if err := studentStore.Ping(ctx); err != nil {
		logrus.Fatalf("Database ping failed: %v", err)
	}
	// Initialize the HTTP handler
	handler := transport.NewHandler(studentService)

	// Start the HTTP server
	if err := handler.Serve(); err != nil {
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

//curl -X GET http://localhost:8080/api/v1/student/1 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTcyMzgwOTUxMX0.T15UN2f28tSu-lagsMIcV8u9Qd2bn_PRKmDCWrzoxkU"

//C:\Users\ADMIN\Desktop\golang-assignment>curl -X POST http://localhost:8080/api/v1/student ^
//More?  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTcyMzgyOTcxOH0.cOmkk5WmPoDRKlWP4iZQtV_I508ST3Xz3C48iJYyqIM" ^
//More? -H "Content-Type: application/json" ^
//More?  -d "{\"id\":\"22\" , \"name\": \"kavya\", \"email\": \"kavya@gmail.com\", \"age\": 21, \"course\": \"Bsc\"}"
