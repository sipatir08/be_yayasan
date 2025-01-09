package donation

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDeleteDonation(t *testing.T) {
	// Membuat router baru di dalam fungsi
	r := mux.NewRouter()
	r.HandleFunc("/donations/{id}", DeleteDonation).Methods("DELETE")

	// Membuat request DELETE ke /donations/{id}
	req, err := http.NewRequest("DELETE", "/donations/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mencatat response dari server
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Mengecek status code response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Mengecek apakah response body sesuai
	expected := `{"message": "Donation deleted successfully"}`
	assert.JSONEq(t, expected, rr.Body.String())
}
