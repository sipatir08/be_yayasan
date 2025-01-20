package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"product/database"
	"product/model/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Register is used to register a new user
func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var existingUser user.User
	err := database.DB.QueryRow("SELECT id FROM users WHERE email=$1", creds.Email).Scan(&existingUser.ID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existingUser.ID != 0 {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error Password", http.StatusInternalServerError)
		return
	}

	// Simpan pengguna baru ke database dengan status belum diverifikasi
	_, err = database.DB.Exec(
		"INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)",
		creds.Name, creds.Email, hashedPassword, "donatur",
	)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Internal server error Insert", http.StatusInternalServerError)
		return
	}

	// Generate verification token
	verificationToken := generateVerificationToken(creds.Email)

	// Kirim email verifikasi
	err = sendVerificationEmail(creds.Email, verificationToken)
	if err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "User registered successfully. Please verify your email.",
	}
	json.NewEncoder(w).Encode(response)
}

// generateVerificationToken creates a token for email verification
func generateVerificationToken(email string) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

// sendVerificationEmail sends an email with the verification link
func sendVerificationEmail(toEmail, verificationToken string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")

	verificationLink := "https://yayasan-three.vercel.app/verify?token=" + verificationToken

	from := mail.NewEmail("Your App", "fathir080604@gmail.com")
	to := mail.NewEmail("User", toEmail)
	subject := "Konfirmasi Registrasi"

	htmlContent := `
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					color: #333;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 30px auto;
					padding: 20px;
					background-color: #ffffff;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				}
				.header {
					text-align: center;
					padding-bottom: 20px;
				}
				.header h2 {
					color: #007BFF;
				}
				.content p {
					font-size: 16px;
				}
				.button {
					display: inline-block;
					padding: 10px 20px;
					background-color: #28a745;
					color: #fff;
					border-radius: 5px;
					text-decoration: none;
					font-size: 16px;
				}
				.footer {
					margin-top: 20px;
					font-size: 14px;
					color: #888;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>Terima kasih telah mendaftar !</h2>
				</div>
				<div class="content">
					<p>Hi,</p>
					<p>Terima kasih telah mendaftar ! Untuk melanjutkan, silakan klik tombol di bawah ini untuk memverifikasi akun Anda:</p>
					<p><a href="` + verificationLink + `" class="button">Verifikasi Akun</a></p>
					<p>Jika Anda tidak mendaftar di website kami, abaikan email ini.</p>
				</div>
				<div class="footer">
					<p>&copy; 2024 OurApp. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	` // The same HTML content as in your original template

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	log.Printf("Email sent successfully: %d", response.StatusCode)
	return nil
}

// VerifyEmail handles email verification
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE users SET is_verified=$1 WHERE email=$2", true, claims.Email)
	if err != nil {
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Email verified successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user user.User
	// Mengambil data user termasuk role dari database
	err := database.DB.QueryRow("SELECT id, email, password, is_verified, role FROM users WHERE email=$1", creds.Email).Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !user.IsVerified {
		http.Error(w, "Please verify your email before logging in", http.StatusForbidden)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Respon JSON yang berisi pesan, token, data pengguna, dan role
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"user": map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"isVerified": user.IsVerified,
			"role":      user.Role,  // Menambahkan role pengguna
		},
	}
	json.NewEncoder(w).Encode(response)
}

// ValidateToken validates the JWT token
func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// JWTAuth is a middleware for JWT authentication
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		sttArr := strings.Split(bearerToken, " ")
		if len(sttArr) == 2 {
			isValid, _ := ValidateToken(sttArr[1])
			if isValid {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
