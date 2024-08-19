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

// Run - sets up our application
func Run() error {

	log.Info("Setting up our application")
	// Set up logging format to JSON and output only to a file
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

//curl -X GET http://localhost:8080/getStudent/1 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTcyMzgwOTUxMX0.T15UN2f28tSu-lagsMIcV8u9Qd2bn_PRKmDCWrzoxkU"

//C:\Users\ADMIN\Desktop\golang-assignment>curl -X POST http://localhost:8080/addStudent ^
//More?  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTcyMzgyOTcxOH0.cOmkk5WmPoDRKlWP4iZQtV_I508ST3Xz3C48iJYyqIM" ^
//More? -H "Content-Type: application/json" ^
//More?  -d "{\"id\":\"22\" , \"name\": \"kavya\", \"email\": \"kavya@gmail.com\", \"age\": 21, \"course\": \"Bsc\"}"

//curl -X DELETE http://localhost:8080/deleteStudent/1 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTcyMzg3MDcxNn0.MfkdH5hPGeTZAyH2ILowlNnFgJaXh45DVDK8RP0j6XU"

//curl -X PUT http://localhost:8080/updateStudent/22 ^
//More? -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTcyMzg3OTI3OH0.O47ozTzNTgxjfKC1H-haPGY2Jkr1pyIo_gISM5_cAuY" ^
//More? -H "Content-Type: application/json" ^
//More? -d "{\"id\": \"22\",\"name\": \"deepa\", \"email\": \"deepa@gmail.com\",\"age\": 25,\"course\": \"MBA\"}"

//curl -X POST http://localhost:8080/login  -H "Content-Type: application/json" -d "{\"user_id\": \"user123\", \"password\": \"password\"}"
