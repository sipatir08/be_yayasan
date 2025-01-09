package donation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"product/database"
	"product/model/donation"
    // "product/model/user"
	"github.com/gorilla/mux"
)

// GetAllDonations retrieves all donations.
func GetAllDonations(w http.ResponseWriter, r *http.Request) {
    var donations []donation.Donation

    query := `SELECT id, nama_donatur, kategori, jumlah, tanggal, type FROM donations`
    rows, err := database.DB.Query(query)
    if err != nil {
        http.Error(w, "Error fetching donations", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var donation donation.Donation
        if err := rows.Scan(&donation.ID, &donation.Nama_donatur, &donation.Kategori, &donation.Jumlah, &donation.Tanggal, &donation.Type); err != nil {
            http.Error(w, "Error reading donation data", http.StatusInternalServerError)
            return
        }
        donations = append(donations, donation)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(donations)
}

// GetDonationByID retrieves a donation by its ID.
func GetDonationByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"]) // Parsing ID dari URL
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var donation donation.Donation
    query := `SELECT id, nama_donatur, kategori, jumlah, tanggal, type FROM donations WHERE id = $1`
    err = database.DB.QueryRow(query, id).Scan(
        &donation.ID, 
        &donation.Nama_donatur, 
        &donation.Kategori, 
        &donation.Jumlah, 
        &donation.Tanggal, 
        &donation.Type,
    )
    
    if err != nil {
        // Mengirimkan respon error dengan status NotFound jika data tidak ditemukan
        http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
        return
    }

    // Mengatur header response untuk konten JSON
    w.Header().Set("Content-Type", "application/json")
    // Mengencode struct donation menjadi JSON dan mengirimkan ke client
    if err := json.NewEncoder(w).Encode(donation); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}

// AddDonation adds a new donation.
func AddDonation(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    err := r.ParseMultipartForm(10 << 20) // Maksimal 10MB untuk file
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Ambil data dari form
    donation := donation.Donation{
        Nama_donatur: r.FormValue("nama_donatur"),
        Kategori:    r.FormValue("kategori"),
        Type:        r.FormValue("type"),
    }

    // Parse jumlah
    jumlah, err := strconv.Atoi(r.FormValue("jumlah"))
    if err != nil {
        http.Error(w, "Invalid amount", http.StatusBadRequest)
        return
    }
    donation.Jumlah = jumlah

    // Ambil file bukti pembayaran
    paymentProof, _, err := r.FormFile("paymentProof")
    if err != nil {
        http.Error(w, "File not found", http.StatusBadRequest)
        return
    }
    defer paymentProof.Close()

    // Simpan bukti pembayaran atau lakukan proses lainnya

    // Validasi data
    if donation.Nama_donatur == "" || donation.Kategori == "" || donation.Jumlah <= 0 || donation.Type == "" {
        http.Error(w, `{"error": "Incomplete donation data"}`, http.StatusBadRequest)
        return
    }

    // Insert donation data into the database
    var id int
    query := `INSERT INTO donations (nama_donatur, kategori, jumlah, type, tanggal) 
              VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
    err = database.DB.QueryRow(query, donation.Nama_donatur, donation.Kategori, donation.Jumlah, donation.Type).Scan(&id)
    if err != nil {
        log.Println("Error adding donation to database:", err)
        http.Error(w, `{"error": "Failed to add donation"}`, http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Donation added successfully", "id": strconv.Itoa(id)})
}

// UpdateDonation updates an existing donation by its ID.
func UpdateDonation(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari parameter URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Decode body JSON ke struct donation
	var donation struct {
		Kategori string  `json:"kategori"`
		Jumlah   float64 `json:"jumlah"`
		Type     string  `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Log data untuk memeriksa input
	fmt.Printf("Received donation update: %+v\n", donation)

	// Validasi input
	if donation.Kategori == "" || donation.Jumlah <= 0 || donation.Type == "" {
		http.Error(w, "Invalid data, please check input fields", http.StatusBadRequest)
		return
	}

	// Update hanya kolom yang dapat diubah
	query := `
		UPDATE donations 
		SET kategori = $1, jumlah = $2, type = $3, tanggal = NOW()
		WHERE id = $4
	`
	_, err = database.DB.Exec(query, donation.Kategori, donation.Jumlah, donation.Type, id)
	if err != nil {
		http.Error(w, "Error updating donation", http.StatusInternalServerError)
		return
	}

	// Kirimkan respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Donation updated successfully"})
}

// DeleteDonation deletes a donation by its ID.
func DeleteDonation(w http.ResponseWriter, r *http.Request) {
    // Ambil ID dari parameter URL
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Log ID untuk memastikan ID yang diterima sesuai
    fmt.Printf("Deleting donation with ID: %d\n", id)

    // Query untuk menghapus donasi
    query := `DELETE FROM donations WHERE id = $1`
    result, err := database.DB.Exec(query, id)
    if err != nil {
        http.Error(w, "Error deleting donation", http.StatusInternalServerError)
        fmt.Printf("Error deleting donation: %v\n", err)
        return
    }

    // Periksa apakah ada baris yang dihapus
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Error checking affected rows", http.StatusInternalServerError)
        fmt.Printf("Error checking affected rows: %v\n", err)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Donation not found", http.StatusNotFound)
        return
    }

    // Kirimkan respons sukses
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Donation deleted successfully"})
}


