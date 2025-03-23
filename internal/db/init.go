package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDBAndMigrate() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConn, err := connectWhithRetry(dsn, 5, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error connecting db: %w", err)
	}

	if err := runMigrations(dsn); err != nil {
		_ = dbConn.Close()
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	log.Println("Migrations successfully applied!")
	return dbConn, nil
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("error initializing migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}

func connectWhithRetry(dsn string, maxAttempts int, delay time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 1; i <= maxAttempts; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			return db, nil
		}
		log.Printf("[DB] Attempt %d/%d failed to connect: %v", i, maxAttempts, err)

		time.Sleep(delay)
	}
	return nil, fmt.Errorf("cannot connect to db after %d attempts: %w", maxAttempts, err)
}
