package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Pool *sql.DB
}

// NewDatabase initializes a new Database instance.
func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Init() error {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	d.Pool, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := d.Pool.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err := d.createTables(); err != nil {
		return fmt.Errorf("failed to create database tables: %v", err)
	}

	log.Println("Connected to the database")
	return nil
}
