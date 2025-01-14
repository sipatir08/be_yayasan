package donation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"product/model/donation"


	"github.com/stretchr/testify/assert"
)

// Setup function to initialize the database connection before tests
func setup() {
	// Initialize your database here
	// For example: database.InitDB()
}

// Test for GetAllDonations
func TestGetAllDonations(t *testing.T) {
	// Mock data
	mockDonations := []donation.Donation{
		{
			ID:           1,
			Nama_donatur: "Donor A",
			Kategori:     "Pendidikan",
			Jumlah:       1000,
			Tanggal:      time.Now(),
			Type:         "Transfer",
		},
		{
			ID:           2,
			Nama_donatur: "Donor B",
			Kategori:     "Kesehatan",
			Jumlah:       2000,
			Tanggal:      time.Now(),
			Type:         "Tunai",
		},
	}

	// Mock handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockDonations)
	})

	// Create HTTP request
	req, err := http.NewRequest("GET", "/donations", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode response body
	var donations []donation.Donation
	err = json.NewDecoder(rr.Body).Decode(&donations)
	if err != nil {
		t.Fatal(err)
	}

	// Assert response data
	assert.Equal(t, len(mockDonations), len(donations))
	for i, d := range donations {
		assert.Equal(t, mockDonations[i].ID, d.ID)
		assert.Equal(t, mockDonations[i].Nama_donatur, d.Nama_donatur)
		assert.Equal(t, mockDonations[i].Kategori, d.Kategori)
		assert.Equal(t, mockDonations[i].Jumlah, d.Jumlah)
		assert.WithinDuration(t, mockDonations[i].Tanggal, d.Tanggal, time.Second) // Periksa waktu dengan toleransi
		assert.Equal(t, mockDonations[i].Type, d.Type)
	}
}

// Test for GetDonationByID
func TestGetDonationByID(t *testing.T) {
	setup()

	// Assuming donation with ID 1 exists in the database
	req, err := http.NewRequest("GET", "/donations/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetDonationByID)

	handler.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert response body
	var donation donation.Donation
	err = json.NewDecoder(rr.Body).Decode(&donation)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the donation details are correct
	assert.Equal(t, 1, donation.ID)
}

// Test for AddDonation
func TestAddDonation(t *testing.T) {
	setup()

	donationData := map[string]interface{}{
		"nama_donatur": "John Doe",
		"kategori":     "Education",
		"jumlah":       1000,
		"type":         "Cash",
	}

	body, err := json.Marshal(donationData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/donations", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddDonation)

	handler.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that donation was added successfully
	assert.Equal(t, "Donation added successfully", response["message"])
}

// Test for UpdateDonation
func TestUpdateDonation(t *testing.T) {
	setup()

	// Assuming donation with ID 1 exists in the database
	donationData := map[string]interface{}{
		"kategori": "Health",
		"jumlah":   2000,
		"type":     "Goods",
	}

	body, err := json.Marshal(donationData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/donations/1", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateDonation)

	handler.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that donation was updated successfully
	assert.Equal(t, "Donation updated successfully", response["message"])
}

// Test for DeleteDonation
func TestDeleteDonation(t *testing.T) {
	setup()

	// Assuming donation with ID 1 exists in the database
	req, err := http.NewRequest("DELETE", "/donations/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteDonation)

	handler.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that donation was deleted successfully
	assert.Equal(t, "Donation deleted successfully", response["message"])
}
