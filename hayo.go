package main

import (
	"product/database"
	"net/http"
	"net/http/httptest"
	"testing"

	"product/controller/auth"
	"product/controller/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMainRoutes(t *testing.T) {
	// Initialize the database
	database.InitDB() // Pastikan database diinisialisasi sebelum pengujian

	// Initialize the router (same as in main.go)
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/users", user.GetUsers).Methods("GET")
	router.HandleFunc("/users", user.AddUser).Methods("POST")
	router.HandleFunc("/users/{id}", user.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", user.DeleteUser).Methods("DELETE")

	// Auth routes
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/regis", auth.Register).Methods("POST")
	router.HandleFunc("/verify", auth.VerifyEmail).Methods("GET")

	// Test GET /users route
	t.Run("GET /users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert that the response code is 200 OK
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test POST /regis route (assuming you have a mock Register handler)
	t.Run("POST /regis", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/regis", nil) // Add request body here
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert that the response code is 200 OK
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
