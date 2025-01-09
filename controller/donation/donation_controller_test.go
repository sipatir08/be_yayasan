package donation

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"product/database" // Pastikan path ini sesuai dengan struktur Anda
)

func TestAddDonation(t *testing.T) {
	// Memanggil fungsi InitDB untuk menginisialisasi koneksi database
	database.InitDB() // Tidak perlu menangkap nilai kembalian

	// Membuat router baru di dalam fungsi
	r := mux.NewRouter()
	r.HandleFunc("/donations", AddDonation).Methods("POST")

	// Data yang akan dikirimkan dalam form
	form := `--boundary
Content-Disposition: form-data; name="nama_donatur"

John Doe
--boundary
Content-Disposition: form-data; name="kategori"

Kebutuhan Panti
--boundary
Content-Disposition: form-data; name="type"

Uang
--boundary
Content-Disposition: form-data; name="jumlah"

100000
--boundary
Content-Disposition: form-data; name="paymentProof"; filename="proof.jpg"
Content-Type: image/jpeg

<image data>
--boundary--`

	// Membuat request POST ke /donations
	req, err := http.NewRequest("POST", "/donations", bytes.NewBufferString(form))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

	// Mencatat response dari server
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Mengecek status code response
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Mengecek apakah response body sesuai
	expected := `{"message": "Donation added successfully", "id": "1"}`
	assert.JSONEq(t, expected, rr.Body.String())
}
