# Hexya NG

This project is a proof-of-concept for a microservices framework inspired by Hexya and Odoo, built with modern technologies.

## Project Structure

The project is organized into the following directories:

- `api-gateway/`: A reverse proxy that routes requests to the appropriate microservice.
- `services/`: Contains the individual microservices.
  - `courses-service/`: Manages courses.
  - `sessions-service/`: Manages course sessions.
  - `attendees-service/`: Manages session attendees.

## Running the Application

To run the application, you need to start each service and the API gateway.

1.  **Start the services:**
    Open a terminal and run each service in the background:
    ```bash
    go run services/courses-service/main.go &
    go run services/sessions-service/main.go &
    go run services/attendees-service/main.go &
    ```

2.  **Start the API Gateway:**
    In the same terminal, run the API gateway:
    ```bash
    go run api-gateway/main.go &
    ```

The services will be running on the following ports:
- `courses-service`: `:8081`
- `sessions-service`: `:8082`
- `attendees-service`: `:8083`
- `api-gateway`: `:8080`

## API Endpoints

All endpoints are accessed through the API gateway on port `:8080`.

### Courses Service

- `GET /api/courses`: List all courses.
- `POST /api/courses`: Create a new course.
  - Body: `{"name": "string", "description": "string"}`
- `GET /api/courses/{id}`: Get a course by ID.
- `PUT /api/courses/{id}`: Update a course.
  - Body: `{"name": "string", "description": "string"}`
- `DELETE /api/courses/{id}`: Delete a course.

### Sessions Service

- `GET /api/sessions`: List all sessions.
- `POST /api/sessions`: Create a new session.
  - Body: `{"name": "string", "course_id": int, "start_date": "string", "duration": int, "seats": int}`
- `GET /api/sessions/{id}`: Get a session by ID.
- `PUT /api/sessions/{id}`: Update a session.
  - Body: `{"name": "string", "course_id": int, "start_date": "string", "duration": int, "seats": int}`
- `DELETE /api/sessions/{id}`: Delete a session.
- `GET /api/courses/{id}/sessions`: List all sessions for a course.

### Attendees Service

- `GET /api/attendees`: List all attendees.
- `POST /api/attendees`: Create a new attendee.
  - Body: `{"name": "string", "session_id": int}`
- `GET /api/attendees/{id}`: Get an attendee by ID.
- `PUT /api/attendees/{id}`: Update an attendee.
  - Body: `{"name": "string", "session_id": int}`
- `DELETE /api/attendees/{id}`: Delete an attendee.
- `GET /api/sessions/{id}/attendees`: List all attendees for a session.
