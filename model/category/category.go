package model

import "time"

// Struktur Data yang sesuai dengan kolom-kolom di database
type Category struct {
	ID           int    `json:"id"`
	Nama_donatur string `json:"nama_donatur"`
	Kategori     string `json:"kategori"`
	Jumlah       int    `json:"jumlah"`
	Tanggal      time.Time `json:"tanggal"`
	Type         string `json:"type"`
}
