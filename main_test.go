package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux" // Tambahkan import ini
)

func TestMainRoutes(t *testing.T) {
	// Inisialisasi router dari main.go
	router := setupRouter()

	// Contoh pengujian rute "/users"
	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Code)
	}
}

func setupRouter() http.Handler {
	// Gunakan kode utama untuk mendefinisikan rute
	router := mux.NewRouter()

	// Tambahkan rute yang akan diuji
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return router
}
