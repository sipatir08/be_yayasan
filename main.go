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
	router.HandleFunc("/category/{category}", controller.GetDataByCategory).Methods("GET")

	// // Admin category routes
	// router.HandleFunc("/admin-category", categoryAdmin.GetCategoryAdmins).Methods("GET")
	// router.HandleFunc("/admin-category", auth.JWTAuth(categoryAdmin.AddCategoryAdmin)).Methods("POST")
	// router.HandleFunc("/admin-category/{id}", auth.JWTAuth(categoryAdmin.UpdateCategoryAdmin)).Methods("PUT")
	// router.HandleFunc("/admin-category/{id}", auth.JWTAuth(categoryAdmin.DeleteCategoryAdmin)).Methods("DELETE")

	// Donation routes
	// router.HandleFunc("/donations", donation.GetDonations).Methods("GET")
	// Donation routes
	router.HandleFunc("/donations", donation.GetAllDonations).Methods("GET")
	router.HandleFunc("/donations/{id}", donation.GetDonationByID).Methods("GET")
	router.HandleFunc("/donations", auth.JWTAuth(donation.AddDonation)).Methods("POST") // JWT middleware diterapkan
	router.HandleFunc("/donations/{id}", donation.UpdateDonation).Methods("PUT")        
	router.HandleFunc("/donations/{id}", donation.DeleteDonation).Methods("DELETE")     


	// CORS settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:5502", "https://yayasann.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Apply CORS to the router
	handler := c.Handler(router)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
