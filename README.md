# golang-assignment

## The directory structure and the role of each file:

1. cmd 
    * (cmd/main.go): Entry point of the application.
    * .env : This file has the environment variables required by the application.
    * app.log : This file stores events, errors, and other messages that are logged by the application.


2. config (config/config.go): This file will manage database configuration settings.

3. internal/student
    * (internal/student/student.go): This will handle student-related logic and data models.
    * (internal/student/login.go): This will authenticate the user and calls a method to generate JWT token.

4. internal/database 
    * (internal/database/student.go and internal/database/database.go): These files will manage database operations and connections.

5. internal/transport
    * (internal/transport/auth.go): This file handles JWT authentication.
    * (internal/transport/handler.go) : This file sets up and manages the HTTP server, routing, and middleware for handling student-related API requests, including CORS, logging, and authentication.
    * (internal/transport/login.go): This file handles user login by validating credentials, authenticating the user, and generating a JWT token for successful logins.
    * (internal/transport/middleware.go): This file defines middleware functions for JSON response formatting, logging, request timeouts, get userID and CORS handling in the application.
    * (internal/transport/srudent.go): This file implements HTTP handlers for managing students, including creating, retrieving, updating, and deleting student records, with validation, JWT authentication, and logging.

6. utils 
    * (utils/jwt.go): Utility functions for JWT token generation.
    * (utils/utils.go): Utility functions for extracting userID and token.
    
* Only admin who is doing the CRUD operations is logging in to the application so there is no entry of login credentials into db, hence I have not written a login.go file in the database package.