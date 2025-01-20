package handler

import (
	"log"
	"net/http"

	"product/controller/auth"
	"product/controller/category"
	"product/controller/donation"
	"product/controller/user"
	"product/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func init() {
	// Inisialisasi database saat server dimulai
	database.InitDB()
}

// Handler function yang digunakan oleh Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Inisialisasi router dengan middleware CORS
	router := SetupRoutes()
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5501", "https://yayasann-sipatir08s-projects.vercel.app/"}, // Disesuaikan
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	// Lanjutkan permintaan ke router
	corsHandler.ServeHTTP(w, r)
}

// Fungsi untuk mengatur rute dan middleware
func SetupRoutes() *mux.Router {
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

	// Category routes
	router.HandleFunc("/category/{category}", category.GetDataByCategory).Methods("GET")

	// Donation routes
	router.HandleFunc("/donations", donation.GetAllDonations).Methods("GET")
	router.HandleFunc("/donations/{id}", donation.GetDonationByID).Methods("GET")
	router.HandleFunc("/donations", auth.JWTAuth(donation.AddDonation)).Methods("POST")
	router.HandleFunc("/donations/{id}", donation.UpdateDonation).Methods("PUT")
	router.HandleFunc("/donations/{id}", donation.DeleteDonation).Methods("DELETE")

	return router
}

// Main function untuk menjalankan server lokal (opsional, tidak diperlukan untuk Vercel)
func main() {
	database.InitDB()

	http.HandleFunc("/", Handler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
