package dbHelper

import (
	"database/sql"
	"fmt"
)

func MigrateDb(db *sql.DB) {
    query := `
		CREATE TABLE IF NOT EXISTS persons (
			id SERIAL PRIMARY KEY,
			phone_number VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
            location VARCHAR(100) NOT NULL
		)`
	_, err := db.Query(query)

    if err != nil {
		fmt.Println("Failed to create table:", err.Error())
	}
}
