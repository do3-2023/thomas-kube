package dbHelper

import (
	"database/sql"
	"log"
)

func PopulateDb(db *sql.DB){

	// create the person table
	createTable := `CREATE TABLE IF NOT EXISTS person (
		id SERIAL NOT NULL PRIMARY KEY,
		last_name TEXT UNIQUE NOT NULL,
		phone_number TEXT NOT NULL,
		location TEXT NOT NULL
	);`

	_, err := db.Exec(createTable)
	if err != nil {
		log.Println("Error creating table:", err)
	}

	sqlStatementInsert := `INSERT INTO person (last_name, phone_number, location)
	VALUES ('John', '0702030405', 'Marseille'),
       ('Doe', '0603040506', 'Montpellier');
	`
	_, errQuery := db.Exec(sqlStatementInsert)
	if errQuery != nil {
		log.Println("Error executing SQL query:", errQuery)
		return
	}

	println("Table setup done")

}
