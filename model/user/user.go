package user

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
	Role       string `json:"role"`
}