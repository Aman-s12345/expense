package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DB_HOST")
	}
	if dbURL == "" {
		dbURL = "postgresql://postgres:password@localhost:5432/expense_tracker?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrationFile := "db/migrations/migration_init_001.sql"
	if len(os.Args) > 1 {
		migrationFile = os.Args[1]
	}

	sqlBytes, err := os.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v", migrationFile, err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	fmt.Println("Migration completed successfully")
}