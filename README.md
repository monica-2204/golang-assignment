# golang-assignment

## The directory structure and the role of each file:

1. Config (config/config.go): This file will manage database configuration settings.
2. Internal (internal/student/student.go): This will handle student-related logic and data models.
3. Database (internal/database/student.go and internal/database/database.go): These files will manage database operations and connections.
4. Transport (transport/auth.go, transport/student.go, transport/handler.go, transport/middleware.go): This layer will handle HTTP requests,       responses, and middleware for authentication.
5. Utils (utils/jwt.go, utils/logging.go): Utility functions for JWT token generation and logging.
6. Main (main.go): Entry point of the application.