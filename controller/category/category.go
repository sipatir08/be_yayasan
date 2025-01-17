package category

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"product/database" // Pastikan path ini sesuai dengan struktur proyek Anda
	model "product/model/category"

	"github.com/gorilla/mux"
)

// Fungsi untuk mengambil data berdasarkan kategori
func GetDataByCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    category := vars["category"] // Mengambil kategori dari URL

    log.Println("Received category:", category)

    // Mengambil data berdasarkan kategori
    donations, err := GetDataByCategoryFromDB(category)
    if err != nil {
        http.Error(w, "No data found or error fetching data", http.StatusInternalServerError)
        return
    }

    // Mengirimkan data sebagai JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(donations)
}

// Fungsi untuk mengambil data dari database berdasarkan kategori
func GetDataByCategoryFromDB(category string) ([]model.Category, error) {
	var data []model.Category

	// Query untuk mengambil data berdasarkan kategori
	query := "SELECT id, nama_donatur, kategori, jumlah, tanggal, type FROM donations WHERE kategori ILIKE $1"
	log.Println("Executing query:", query, "with category:", category)
	rows, err := database.DB.Query(query, category)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Memasukkan hasil query ke dalam slice data
	for rows.Next() {
		var d model.Category
		err := rows.Scan(&d.ID, &d.Nama_donatur, &d.Kategori, &d.Jumlah, &d.Tanggal, &d.Type)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		data = append(data, d)
	}

	// Memastikan ada data yang ditemukan
	if len(data) == 0 {
		return nil, errors.New("no data found")
	}

	return data, nil
}
