package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

//go:embed static/index.html
var staticFiles embed.FS

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type Car struct {
	ID          int     `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Year        int     `json:"year"`
	Price       float64 `json:"price"`
	FuelType    string  `json:"fuelType"`
	Horsepower  int     `json:"horsepower"`
	Mileage     int     `json:"mileage"`
	Image       string  `json:"image"`
	Featured    bool    `json:"featured"`
	Description string  `json:"description"`
}

type DealershipStats struct {
	TotalCars     int     `json:"totalCars"`
	AvgPrice      float64 `json:"avgPrice"`
	TotalValue    float64 `json:"totalValue"`
	BrandsCount   int     `json:"brandsCount"`
}

var cars = []Car{
	{
		ID:          1,
		Brand:       "Porsche",
		Model:       "911 Turbo S",
		Year:        2024,
		Price:       225000,
		FuelType:    "Gasoline",
		Horsepower:  640,
		Mileage:     150,
		Image:       "https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e?w=400&h=300&fit=crop",
		Featured:    true,
		Description: "The ultimate daily-drivable supercar with breathtaking performance.",
	},
	{
		ID:          2,
		Brand:       "Mercedes-Benz",
		Model:       "AMG GT 63 S",
		Year:        2024,
		Price:       175000,
		FuelType:    "Gasoline",
		Horsepower:  630,
		Mileage:     200,
		Image:       "https://images.unsplash.com/photo-1605559424843-9e4c228bf1c2?w=400&h=300&fit=crop",
		Featured:    true,
		Description: "Four-door supercar blending luxury with track-ready performance.",
	},
	{
		ID:          3,
		Brand:       "Tesla",
		Model:       "Model S Plaid",
		Year:        2024,
		Price:       89990,
		FuelType:    "Electric",
		Horsepower:  1020,
		Mileage:     50,
		Image:       "https://images.unsplash.com/photo-1617788138017-80ad40651399?w=400&h=300&fit=crop",
		Featured:    true,
		Description: "Mind-blowing acceleration with cutting-edge electric technology.",
	},
	{
		ID:          4,
		Brand:       "BMW",
		Model:       "M8 Competition",
		Year:        2023,
		Price:       140000,
		FuelType:    "Gasoline",
		Horsepower:  617,
		Mileage:     5000,
		Image:       "https://images.unsplash.com/photo-1617531653332-bd46c24f2068?w=400&h=300&fit=crop",
		Featured:    false,
		Description: "Grand tourer with brutal power and elegant styling.",
	},
	{
		ID:          5,
		Brand:       "Audi",
		Model:       "RS e-tron GT",
		Year:        2024,
		Price:       147000,
		FuelType:    "Electric",
		Horsepower:  637,
		Mileage:     100,
		Image:       "https://images.unsplash.com/photo-1614200187524-dc4b892acf16?w=400&h=300&fit=crop",
		Featured:    false,
		Description: "Electrifying performance meets Audi's signature sophistication.",
	},
	{
		ID:          6,
		Brand:       "Lamborghini",
		Model:       "Huracán STO",
		Year:        2023,
		Price:       331000,
		FuelType:    "Gasoline",
		Horsepower:  631,
		Mileage:     1200,
		Image:       "https://images.unsplash.com/photo-1621135802920-133df287f89c?w=400&h=300&fit=crop",
		Featured:    true,
		Description: "Race-inspired V10 monster derived from Lamborghini's motorsport DNA.",
	},
	{
		ID:          7,
		Brand:       "Ferrari",
		Model:       "Roma",
		Year:        2024,
		Price:       250000,
		FuelType:    "Gasoline",
		Horsepower:  612,
		Mileage:     80,
		Image:       "https://images.unsplash.com/photo-1592198084033-aade902d79a3?w=400&h=300&fit=crop",
		Featured:    false,
		Description: "Timeless Italian elegance with exhilarating V8 performance.",
	},
	{
		ID:          8,
		Brand:       "Lucid",
		Model:       "Air Sapphire",
		Year:        2024,
		Price:       249000,
		FuelType:    "Electric",
		Horsepower:  1200,
		Mileage:     30,
		Image:       "https://images.unsplash.com/photo-1617469165786-8007eda3caa7?w=400&h=300&fit=crop",
		Featured:    false,
		Description: "The new standard for electric luxury sedans with insane power.",
	},
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "ok",
		Version:   "2.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func carsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Support query params for filtering
	brand := r.URL.Query().Get("brand")
	featured := r.URL.Query().Get("featured") == "true"
	
	filteredCars := cars
	if brand != "" {
		temp := []Car{}
		for _, car := range cars {
			if car.Brand == brand {
				temp = append(temp, car)
			}
		}
		filteredCars = temp
	}
	
	if featured {
		temp := []Car{}
		for _, car := range filteredCars {
			if car.Featured {
				temp = append(temp, car)
			}
		}
		filteredCars = temp
	}
	
	json.NewEncoder(w).Encode(filteredCars)
}

func carDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/car/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}
	
	for _, car := range cars {
		if car.ID == id {
			json.NewEncoder(w).Encode(car)
			return
		}
	}
	http.Error(w, "Car not found", http.StatusNotFound)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	brandsMap := make(map[string]bool)
	var totalValue float64
	
	for _, car := range cars {
		brandsMap[car.Brand] = true
		totalValue += car.Price
	}
	
	stats := DealershipStats{
		TotalCars:   len(cars),
		AvgPrice:    totalValue / float64(len(cars)),
		TotalValue:  totalValue,
		BrandsCount: len(brandsMap),
	}
	
	json.NewEncoder(w).Encode(stats)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content, _ := staticFiles.ReadFile("static/index.html")
	w.Write(content)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/cars", carsHandler)
	http.HandleFunc("/car/", carDetailHandler)
	http.HandleFunc("/stats", statsHandler)
	
	log.Println(" LuxDrive Dealership API v2.0 starting on :8080")
	log.Println(" Modern luxury car marketplace ready at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}