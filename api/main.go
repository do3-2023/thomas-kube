package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thomas-mauran/city_api/dbHelper"
	city "github.com/thomas-mauran/city_api/struct"
	"github.com/thomas-mauran/city_api/utils"
)

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "path", "status"},
)

func incrementRequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.WithLabelValues(r.Method, r.URL.Path, http.StatusText(http.StatusOK)).Inc()
		next.ServeHTTP(w, r)
	})
}

func main() {
	cityAPIAddr := os.Getenv("CITY_API_ADDR")
	cityAPIPort := os.Getenv("CITY_API_PORT")
	cityAPIDBUrl := os.Getenv("CITY_API_DB_URL")
	cityAPIDbUser := os.Getenv("CITY_API_DB_USER")

	fmt.Println("cityAPIAddr:", cityAPIAddr)
	fmt.Println("cityAPIPort:", cityAPIPort)

	if cityAPIDBUrl == "" || cityAPIDbUser == ""  {
		log.Fatal("Missing some environment variables")
	}

	// Connect to the database
	db, err := sql.Open("postgres", cityAPIDBUrl)
	if err != nil {
		println("Error opening the database connection:", err)
		panic(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Prometheus metrics
	prometheus.MustRegister(requestCount)

	println("Starting the server on :", cityAPIAddr+":"+cityAPIPort)

	// Populate the db
	dbHelper.PopulateDb(db)
	println("test")

	r.Route("/metrics", func(r chi.Router) {
		r.Use(incrementRequestCount)
		r.Handle("/", promhttp.Handler())
	})

	// Health check
	r.Get("/_health", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			log.Println("Unable to ping the database:", err)
			response := fmt.Sprintf("Unable to ping the database: %v", err)
			utils.Response(w, r, 500, response)
			return
		}
		log.Println("Everything is good!")
		utils.Response(w, r, 204, "Everything is good!")
	})

	// City GET
	r.Get("/city", func(w http.ResponseWriter, r *http.Request) {
		sqlStatement := `SELECT * FROM city`
		rows, err := db.Query(sqlStatement)
		if err != nil {
			log.Println("Error querying the database:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer rows.Close()

		var listOfCities []city.City

		for rows.Next() {
			var city city.City
			if err := rows.Scan(&city.ID, &city.DepartmentCode, &city.InseeCode, &city.ZipCode, &city.Name, &city.Lat, &city.Lon); err != nil {
				log.Println("Error scanning row:", err)
				utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			listOfCities = append(listOfCities, city)
		}
		if err := rows.Err(); err != nil {
			log.Println("Error iterating rows:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		jsonData, err := json.Marshal(listOfCities)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonData); err != nil {
			log.Println("Error writing response:", err)
			// Handle the error. You can choose to log the error, send an appropriate response, or take any other action.
		}
	})

	// City POST
	r.Post("/city", func(w http.ResponseWriter, r *http.Request) {
		var cityObj city.City
		err := json.NewDecoder(r.Body).Decode(&cityObj)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		sqlStatement := `INSERT INTO city (department_code, insee_code, zip_code, name, lat, lon) VALUES ($1, $2, $3, $4, $5, $6)`
		_, errQuery := db.Exec(sqlStatement, cityObj.DepartmentCode, cityObj.InseeCode, cityObj.ZipCode, cityObj.Name, cityObj.Lat, cityObj.Lon)
		if errQuery != nil {
			log.Println("Error executing SQL query:", errQuery)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.Response(w, r, http.StatusCreated, "Posted!")
	})

	errServ := http.ListenAndServe(cityAPIAddr + ":" + cityAPIPort, r)
	if errServ != nil {
		log.Fatal(errServ)
	}
}
