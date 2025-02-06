# Golang JWT-Based REST API

## Overview
This is a simple REST API built with Golang that uses JSON Web Tokens (JWT) for authentication and role-based access control. The API allows users to log in and retrieve tasks, while administrators can add new tasks.

## Features
- User authentication via JWT
- Role-based access control (Admin vs User)
- Task management
- RESTful API design

## Endpoints

### 1. User Authentication
#### `POST /api/login`
Authenticates users and returns a JWT.

**Request Body (JSON):**
```json
{
  "username": "admin",
  "password": "adminpass"
}
```

**Response:**
```json
{
  "token": "your-jwt-token"
}
```

**Valid Credentials:**
| Username | Password  | Role  |
|----------|----------|-------|
| admin    | adminpass | admin |
| user     | userpass  | user  |

### 2. Task Management
#### `GET /api/tasks`
Retrieves the list of tasks.
- Available for both **admin** and **user** roles.

**Response:**
```json
[
  { "task": "Learn Golang" },
  { "task": "Build a REST API in Golang" }
]
```

#### `POST /api/tasks`
Adds a new task (Admin only).

**Request Body (JSON):**
```json
{
  "task": "New Task Name"
}
```

**Response:**
```json
{
  "task": "New Task Name"
}
```

#### `PUT /api/tasks` (Not allowed)
#### `DELETE /api/tasks` (Not allowed)
- Only GET and POST methods are supported.

## Setup and Running the Server

1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo.git
   cd your-repo
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Run the server:
   ```sh
   go run main.go
   ```
4. The server starts at `http://localhost:8080`

## Libraries Used
- `net/http` - for handling HTTP requests
- `encoding/json` - for JSON encoding/decoding
- `github.com/dgrijalva/jwt-go` - for JWT authentication

## Notes
- Ensure `secretKey` is kept secure.
- The token expiration is set to **24 hours**.
- Role-based authorization restricts task modifications to **admins only**.



