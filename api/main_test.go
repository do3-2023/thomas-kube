package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	city "github.com/thomas-mauran/city_api/struct"
)

func TestGetCities(t *testing.T) {
	req, err := http.NewRequest("GET", "/city", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCitiesHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	var cities []city.City
	err = json.Unmarshal(rr.Body.Bytes(), &cities)
	if err != nil {
		t.Errorf("Error unmarshaling JSON response: %v", err)
	}

	// Perform assertions on the cities variable as needed
	// ...
}

func TestCreateCity(t *testing.T) {
	city := city.City{
		ID:             1,
		DepartmentCode: "01",
		InseeCode:      "01001",
		ZipCode:        "01000",
		Name:           "Test City",
		Lat:            0.0,
		Lon:            0.0,
	}

	jsonBody, err := json.Marshal(city)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/city", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createCityHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Perform assertions on the response as needed
	// ...
}

func getCitiesHandler(w http.ResponseWriter, r *http.Request) {
	// Mock the database query and return dummy data
	cities := []city.City{
		{ID: 1, DepartmentCode: "01", InseeCode: "01001", ZipCode: "01000", Name: "City 1", Lat: 0.0, Lon: 0.0},
		{ID: 2, DepartmentCode: "02", InseeCode: "02001", ZipCode: "02000", Name: "City 2", Lat: 0.0, Lon: 0.0},
	}

	jsonData, err := json.Marshal(cities)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func createCityHandler(w http.ResponseWriter, r *http.Request) {
	var city city.City
	err := json.NewDecoder(r.Body).Decode(&city)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Mock the database insert operation

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte("Posted!")) ; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	
}
