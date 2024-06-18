package dbHelper

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func MigrateDb(db *sql.DB) {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if (err != nil) {
        log.Fatalf("Could not create DB driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations", 
        "postgres", driver)
    if (err != nil) {
        log.Fatalf("Could not create migrate instance: %v", err)
    }

    err = m.Up()
    if (err != nil && err != migrate.ErrNoChange) {
        log.Fatalf("Could not run migrations: %v", err)
    }

    log.Println("Database migrated successfully")
}
