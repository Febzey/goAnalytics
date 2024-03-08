package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Database struct {
	db *sql.DB
	mu *sync.Mutex
}

type PageAnalyticData struct {
	Event     string
	UserAgent string
	URL       string
	Referrer  string
}

type PageView struct {
	PageURL   string
	IPAddress string
	LastView  time.Time
}

// NewDatabase initializes a new Database instance.
func NewDatabase() *Database {
	return &Database{
		mu: &sync.Mutex{},
	}
}

func (d *Database) Init() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.db != nil {
		return nil // Database already initialized
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	d.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Connected to the database")
	return nil
}

func (d *Database) InsertPageAnalyticData() {

}
