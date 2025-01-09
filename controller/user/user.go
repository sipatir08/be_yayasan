package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"product/database"
	"product/model/user"

	"github.com/gorilla/mux"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (name, email, password, role, is_verified) VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.Role, u.IsVerified)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updatedUser user.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	var queryParts []string
	var args []interface{}

	if updatedUser.Name != "" {
		queryParts = append(queryParts, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Name)
	}
	if updatedUser.Email != "" {
		queryParts = append(queryParts, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Email)
	}
	if updatedUser.Password != "" {
		queryParts = append(queryParts, "password = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Password)
	}
	if updatedUser.Role != "" {
		queryParts = append(queryParts, "role = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.Role)
	}
	if updatedUser.IsVerified {
		queryParts = append(queryParts, "is_verified = $"+strconv.Itoa(len(args)+1))
		args = append(args, updatedUser.IsVerified)
	}

	if len(queryParts) == 0 {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET " + strings.Join(queryParts, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)

	_, err = database.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err := database.DB.Exec(query, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`SELECT id, name, email, password, role, is_verified FROM users`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsVerified); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
