package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	dbURL := "postgres://postgres:root@localhost:5432/grpc-auth?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic("unable to connect to postgres")
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic("unable to close postgres connection")
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic("unable to configure database")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"grpc-auth",
		driver,
	)
	if err != nil {
		panic("Unable to Initialize migrate")
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")

			return
		}

		panic("Unable to Run migration")
	}

	fmt.Println("Migrations ran successfully!")
}
