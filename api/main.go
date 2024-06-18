package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/do3-2023/thomas-kube/dbHelper"
	person "github.com/do3-2023/thomas-kube/struct"
	"github.com/do3-2023/thomas-kube/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	APIAddr := os.Getenv("API_ADDR")
	APIPort := os.Getenv("API_PORT")
	APIDBUrl := os.Getenv("API_DB_URL")
	APIDbUser := os.Getenv("API_DB_USER")

	fmt.Println("APIAddr:", APIAddr)
	fmt.Println("APIPort:", APIPort)
	fmt.Println("APIDbUrl:", APIDBUrl)
	fmt.Println("APIDbUser:", APIDbUser)

	if APIDBUrl == "" || APIDbUser == ""  {
		log.Fatal("Missing some environment variables")
	}

	// Connect to the database
	const maxRetries = 5
	const retryDelay = 5 * time.Second // Adjust the delay time as needed
	
	var db *sql.DB
	
	for retries := 0; retries < maxRetries; retries++ {
		var err error
		db, err = sql.Open("postgres", APIDBUrl)
		if err != nil {
			log.Printf("Error opening the database connection (attempt %d/%d): %v\n", retries+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}
	
		// Check if the database is accessible by pinging it
		err = db.Ping()
		if err != nil {
			log.Printf("Error connecting to the database (attempt %d/%d): %v\n", retries+1, maxRetries, err)
			db.Close()
			time.Sleep(retryDelay)
			continue
		}
	
		// Connection successful, break out of the loop
		log.Println("Database connection established successfully.")
		break
	}
	
	if db == nil {
		log.Fatal("Unable to establish database connection after maximum retries.")
	}

	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Prometheus metrics
	prometheus.MustRegister(requestCount)

	println("Starting the server on :", APIAddr+":"+APIPort)

	// Populate the db
    dbHelper.MigrateDb(db)
	println("test")

	r.Route("/metrics", func(r chi.Router) {
		r.Use(incrementRequestCount)
		r.Handle("/", promhttp.Handler())
	})

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			log.Println("Unable to ping the database:", err)
			response := fmt.Sprintf("Unable to ping the database: %v", err)
			utils.Response(w, r, 500, response)
			return
		}
		log.Println("Everything is good!")
		utils.Response(w, r, 200, "Everything is good!")
	})

	// Person GET
	r.Get("/person", func(w http.ResponseWriter, r *http.Request) {
		sqlStatement := `SELECT * FROM person`
		rows, err := db.Query(sqlStatement)
		if err != nil {
			log.Println("Error querying the database:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer rows.Close()

		var listOfPersons []person.Person

		for rows.Next() {
			var person person.Person
			if err := rows.Scan(&person.ID, &person.LastName, &person.PhoneNumber, &person.Location); err != nil {
				log.Println("Error scanning row:", err)
				utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			listOfPersons = append(listOfPersons, person)
		}
		if err := rows.Err(); err != nil {
			log.Println("Error iterating rows:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		jsonData, err := json.Marshal(listOfPersons)
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

	//Random person
	r.Get("/person/random", func(w http.ResponseWriter, r *http.Request) {
		sqlStatement := `SELECT * FROM person ORDER BY RANDOM() LIMIT 1`
		row := db.QueryRow(sqlStatement)

		var person person.Person
		if err := row.Scan(&person.ID, &person.LastName, &person.PhoneNumber, &person.Location); err != nil {
			log.Println("Error scanning row:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		jsonData, err := json.Marshal(person)
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

	// Person POST
	r.Post("/person", func(w http.ResponseWriter, r *http.Request) {
		var Obj person.Person
		err := json.NewDecoder(r.Body).Decode(&Obj)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		sqlStatement := `INSERT INTO person (last_name, phone_number, location) VALUES ($1, $2, $3)`
		_, errQuery := db.Exec(sqlStatement, Obj.LastName, Obj.PhoneNumber, Obj.Location)
		if errQuery != nil {
			log.Println("Error executing SQL query:", errQuery)
			utils.Response(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.Response(w, r, http.StatusCreated, "Posted!")
	})

	errServ := http.ListenAndServe(":" + APIPort, r)
	if errServ != nil {
		log.Fatal(errServ)
	}
}
