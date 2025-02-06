package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your-secret-key") // Ensure this is declared only once

// Sample data
var tasks = []map[string]string{
	{"task": "Learn Golang"},
	{"task": "Build a REST API in Golang"},
}

func GenerateJWT(username, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		role := claims["role"].(string)
		return username, role, nil
	}
	return "", "", fmt.Errorf("invalid token")
}

func authenticateJWT(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return "", "", false
	}

	tokenString := strings.Split(authHeader, " ")[1]

	username, role, err := ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return "", "", false
	}

	return username, role, true
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds map[string]string
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Assign role based on the username
	var role string
	if creds["username"] == "admin" {
		role = "admin"
	} else {
		role = "user"
	}

	if creds["username"] == "admin" && creds["password"] == "adminpass" ||
		creds["username"] == "user" && creds["password"] == "userpass" {

		token, err := GenerateJWT(creds["username"], role) // Call GenerateJWT with role
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	_, role, valid := authenticateJWT(w, r)
	if !valid {
		return
	}

	if role == "admin" {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
		} else if r.Method == http.MethodPost {
			var newTask map[string]string
			json.NewDecoder(r.Body).Decode(&newTask)
			tasks = append(tasks, newTask)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newTask)
		} else if r.Method == http.MethodPut || r.Method == http.MethodDelete {

			http.Error(w, "Method not allowed for non-admins", http.StatusMethodNotAllowed)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if role == "user" {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
		} else {
			http.Error(w, "Method not allowed for users", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, "Invalid role", http.StatusUnauthorized)
	}
}

func main() {
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/tasks", tasksHandler)

	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
