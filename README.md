# golang-assignment

## Here is a breakdown of the directory structure and the role of each file:

Config (config/config.go): This file will manage database configuration settings.
Internal (internal/student/student.go): This will handle student-related logic and data models.
Database (internal/database/student.go and internal/database/database.go): These files will manage database operations and connections.
Transport (transport/auth.go, transport/student.go, transport/handler.go, transport/middleware.go): This layer will handle HTTP requests, responses, and middleware for authentication.
Utils (utils/jwt.go, utils/logging.go): Utility functions for JWT token generation and logging.
Main (main.go): Entry point of the application.