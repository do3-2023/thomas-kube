package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	person "github.com/do3-2023/thomas-kube/struct"

	"github.com/do3-2023/thomas-kube/dbHelper"
	"github.com/fermyon/spin-go-sdk/variables"
	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/pg"
)

type Database struct {
	db *sql.DB
}
func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		// create connection to DB
		APIDBUrl, _ := variables.Get("db_url")
		APIDbUser, _ := variables.Get("db_username")
		APIDbPassword, _ := variables.Get("db_password")
		APIDbName, _ := variables.Get("db_name")

		if APIDBUrl == "" || APIDbUser == "" {
			log.Fatal("Missing some environment variables")
		}

		
		var connexion = pg.Open(fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
			APIDBUrl,
			APIDbUser,
			APIDbPassword,
			APIDbName),
		)

		var db = &Database{connexion}
		

		fmt.Println("Connected to DB")

		dbHelper.MigrateDb(db.db)


		router := spinhttp.NewRouter()
		router.GET("/persons", db.handleGetPersons)
		router.POST("/persons", db.handlePostPerson) 

		router.ServeHTTP(w, r)
	})
}

func main() {}

func (db *Database) handleGetPersons(w http.ResponseWriter, r *http.Request, ps spinhttp.Params) {
	err := db.db.Ping()
	if err != nil {
		fmt.Println("Can't ping db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("GET REQUEST ON /persons")

	fmt.Println("db", &db)

	rows, err := db.db.Query("SELECT last_name, phone_number, location FROM persons")
	if err != nil {
		fmt.Println("BIG ERROR")
		fmt.Println(db.db.Stats())
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



	var persons []*person.Person

	println("Persons:", rows)

	for rows.Next() {
		var person person.Person
		if err := rows.Scan(&person.LastName, &person.PhoneNumber, &person.Location); err != nil {
			fmt.Println("Error scanning row:", err)
		}
		persons = append(persons, &person)
	}


	// convert to json
	jsonData, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonData))

}

func (db *Database) handlePostPerson(w http.ResponseWriter, r *http.Request, ps spinhttp.Params) {

	// log the r 
	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error reading body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var Obj person.Person

	err = json.Unmarshal(body, &Obj)

	if err != nil {
		fmt.Println("Error unmarshalling JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sqlStatement := `INSERT INTO persons (last_name, phone_number, location) VALUES ($1, $2, $3)`

	_, err = db.db.Exec(sqlStatement, Obj.LastName, Obj.PhoneNumber, Obj.Location)

	if err != nil {
		fmt.Println("Error executing SQL query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Posted!")
}
