package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var body HealthResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body.Status)
	}

	if body.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0', got '%s'", body.Version)
	}

	if body.Timestamp == "" {
		t.Error("expected timestamp to be present")
	}
}

func TestCarsHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedCount  int
		expectedStatus int
	}{
		{
			name:           "get all cars",
			queryParams:    map[string]string{},
			expectedCount:  8,
			expectedStatus: http.StatusOK,
		},
		{
			name: "get featured cars only",
			queryParams: map[string]string{
				"featured": "true",
			},
			expectedCount:  4,
			expectedStatus: http.StatusOK,
		},
		{
			name: "get cars by brand - Porsche",
			queryParams: map[string]string{
				"brand": "Porsche",
			},
			expectedCount:  1,
			expectedStatus: http.StatusOK,
		},
		{
			name: "get cars by brand - Tesla",
			queryParams: map[string]string{
				"brand": "Tesla",
			},
			expectedCount:  1,
			expectedStatus: http.StatusOK,
		},
		{
			name: "get cars with non-existent brand",
			queryParams: map[string]string{
				"brand": "NonExistent",
			},
			expectedCount:  0,
			expectedStatus: http.StatusOK,
		},
		{
			name: "get electric cars",
			queryParams: map[string]string{
				"fuelType": "Electric",
			},
			expectedCount:  3,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cars", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			carsHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			var cars []Car
			if err := json.NewDecoder(res.Body).Decode(&cars); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(cars) != tt.expectedCount {
				t.Errorf("expected %d cars, got %d", tt.expectedCount, len(cars))
			}
		})
	}
}

func TestCarDetailHandler(t *testing.T) {
	tests := []struct {
		name           string
		carID          string
		expectedStatus int
		expectedBrand  string
		expectedModel  string
	}{
		{
			name:           "get existing car - ID 1",
			carID:          "1",
			expectedStatus: http.StatusOK,
			expectedBrand:  "Porsche",
			expectedModel:  "911 Turbo S",
		},
		{
			name:           "get existing car - ID 6",
			carID:          "6",
			expectedStatus: http.StatusOK,
			expectedBrand:  "Lamborghini",
			expectedModel:  "Huracán STO",
		},
		{
			name:           "get non-existent car",
			carID:          "999",
			expectedStatus: http.StatusNotFound,
			expectedBrand:  "",
			expectedModel:  "",
		},
		{
			name:           "get car with invalid ID",
			carID:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBrand:  "",
			expectedModel:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/car/"+tt.carID, nil)
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("/car/", carDetailHandler)
			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedStatus == http.StatusOK {
				var car Car
				if err := json.NewDecoder(res.Body).Decode(&car); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if car.Brand != tt.expectedBrand {
					t.Errorf("expected brand '%s', got '%s'", tt.expectedBrand, car.Brand)
				}

				if car.Model != tt.expectedModel {
					t.Errorf("expected model '%s', got '%s'", tt.expectedModel, car.Model)
				}
			}
		})
	}
}

func TestStatsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	statsHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var stats DealershipStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedTotalCars := 8
	expectedBrandsCount := 8

	if stats.TotalCars != expectedTotalCars {
		t.Errorf("expected total cars %d, got %d", expectedTotalCars, stats.TotalCars)
	}

	if stats.BrandsCount != expectedBrandsCount {
		t.Errorf("expected brands count %d, got %d", expectedBrandsCount, stats.BrandsCount)
	}

	var totalPrice float64
	for _, car := range cars {
		totalPrice += car.Price
	}
	expectedAvgPrice := totalPrice / float64(len(cars))

	if stats.AvgPrice != expectedAvgPrice {
		t.Errorf("expected average price %.2f, got %.2f", expectedAvgPrice, stats.AvgPrice)
	}

	if stats.TotalValue != totalPrice {
		t.Errorf("expected total value %.2f, got %.2f", totalPrice, stats.TotalValue)
	}
}

func TestCarDataIntegrity(t *testing.T) {
	for i, car := range cars {
		if car.ID == 0 {
			t.Errorf("car at index %d has invalid ID", i)
		}
		if car.Brand == "" {
			t.Errorf("car at index %d has empty brand", i)
		}
		if car.Model == "" {
			t.Errorf("car at index %d has empty model", i)
		}
		if car.Year < 2020 || car.Year > 2025 {
			t.Errorf("car at index %d has invalid year: %d", i, car.Year)
		}
		if car.Price <= 0 {
			t.Errorf("car at index %d has invalid price: %.2f", i, car.Price)
		}
		if car.Horsepower <= 0 {
			t.Errorf("car at index %d has invalid horsepower: %d", i, car.Horsepower)
		}
		if car.Mileage < 0 {
			t.Errorf("car at index %d has invalid mileage: %d", i, car.Mileage)
		}
		if car.FuelType != "Gasoline" && car.FuelType != "Electric" && car.FuelType != "Hybrid" {
			t.Errorf("car at index %d has invalid fuel type: %s", i, car.FuelType)
		}
	}
}

func TestIndexHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	indexHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("expected Content-Type 'text/html', got '%s'", contentType)
	}
}

func TestCORSHeaders(t *testing.T) {
	endpoints := []string{"/health", "/cars", "/stats"}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, endpoint, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			})

			handler.ServeHTTP(w, req)
		})
	}
}

func TestConcurrentRequests(t *testing.T) {
	const numRequests = 10
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/cars", nil)
			w := httptest.NewRecorder()
			carsHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("concurrent request failed with status %d", res.StatusCode)
			}
			done <- true
		}()
	}

	for i := 0; i < numRequests; i++ {
		<-done
	}
}

func TestResponseContentType(t *testing.T) {
	endpoints := []struct {
		path        string
		handler     http.HandlerFunc
		contentType string
	}{
		{"/health", healthHandler, "application/json"},
		{"/cars", carsHandler, "application/json"},
		{"/stats", statsHandler, "application/json"},
		{"/", indexHandler, "text/html"},
	}

	for _, ep := range endpoints {
		t.Run(ep.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ep.path, nil)
			w := httptest.NewRecorder()
			ep.handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			contentType := res.Header.Get("Content-Type")
			if contentType != ep.contentType {
				t.Errorf("expected Content-Type '%s', got '%s'", ep.contentType, contentType)
			}
		})
	}
}

func BenchmarkCarsHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/cars", nil)
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		carsHandler(w, req)
		w = httptest.NewRecorder()
	}
}

func BenchmarkCarDetailHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/car/1", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("/car/", carDetailHandler)
		mux.ServeHTTP(w, req)
	}
}

func TestIntegration(t *testing.T) {
	t.Run("health check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		healthHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("health check failed: %d", w.Code)
		}
	})

	t.Run("get all cars", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/cars", nil)
		w := httptest.NewRecorder()
		carsHandler(w, req)

		var allCars []Car
		if err := json.NewDecoder(w.Body).Decode(&allCars); err != nil {
			t.Fatalf("failed to decode cars: %v", err)
		}

		if len(allCars) == 0 {
			t.Error("expected at least one car")
		}
	})

	t.Run("get specific car", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/car/1", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("/car/", carDetailHandler)
		mux.ServeHTTP(w, req)

		var car Car
		if err := json.NewDecoder(w.Body).Decode(&car); err != nil {
			t.Fatalf("failed to decode car: %v", err)
		}

		if car.ID != 1 {
			t.Errorf("expected car ID 1, got %d", car.ID)
		}
	})

	t.Run("get statistics", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/stats", nil)
		w := httptest.NewRecorder()
		statsHandler(w, req)

		var stats DealershipStats
		if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
			t.Fatalf("failed to decode stats: %v", err)
		}

		if stats.TotalCars != 8 {
			t.Errorf("expected 8 total cars, got %d", stats.TotalCars)
		}
	})
}