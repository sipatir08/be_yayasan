package main

import (
	"log"
	"net/http"
	"os"

	"product/controller/auth"
	"product/controller/category"
	"product/controller/donation"
	"product/controller/user"
	"product/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if os.Getenv("VERCEL_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Initialize the database
	database.InitDB()

	// Initialize the router
	router := mux.NewRouter()

	// Setup routes
	SetupRoutes(router)

	// CORS settings
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "http://localhost:3000", "http://127.0.0.1:5502", "https://yayasann.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Start the server
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", corsHandler.Handler(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Fungsi untuk mengatur rute dan middleware
func SetupRoutes(router *mux.Router) {
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
}

// Fungsi handler untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Setup routes
	SetupRoutes(router)

	// CORS settings
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://yayasann-sipatir08s-projects.vercel.app", "http://localhost:8080/donations", "http://localhost:3000", "http://127.0.0.1:5502", "https://yayasann.vercel.app", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Apply CORS middleware
	corsHandler.Handler(router).ServeHTTP(w, r)
}
