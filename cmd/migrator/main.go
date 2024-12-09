package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"go-music-lib/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.Parse()

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	cfg, err := config.MustLoad()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic("Невозможно подключиться к PostgreSQL" + err.Error())
	}

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatalf("Невозможно закрыть соединение с PostgreSQL: %v", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось проверить соединение с базой данных: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Невозможно настроить драйвер базы данных: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		cfg.DBName,
		driver,
	)
	if err != nil {
		panic("Невозможно инициализировать миграцию: %v" + err.Error())
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Нет миграций для применения")
			return
		}
		log.Fatalf("Невозможно выполнить миграцию: %v", err)
	}

	fmt.Println("Миграции успешно применены!")
}
