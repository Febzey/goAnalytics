package database

import (
	"fmt"
	"log"
	"time"
)

// Page represents the pages table.
type Page struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	ViewCount int    `json:"view_count"`
}

// PageView represents the page_views table.
type PageView struct {
	ID        int64     `json:"id"`
	PageID    int64     `json:"page_id"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
}

func (d *Database) createTables() error {
	// Execute SQL statements to create tables
	statements := []string{
		`CREATE TABLE IF NOT EXISTS pages (
			id INT AUTO_INCREMENT PRIMARY KEY,
			url VARCHAR(255) NOT NULL,
			view_count INT DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS page_views (
			id INT AUTO_INCREMENT PRIMARY KEY,
			page_id INT,
			user_agent VARCHAR(255),
			referrer VARCHAR(255),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			ip_address VARCHAR(45),
			FOREIGN KEY (page_id) REFERENCES pages(id)
		);`,
	}

	for _, statement := range statements {
		_, err := d.db.Exec(statement)
		if err != nil {
			return fmt.Errorf("failed to create tables: %v", err)
		}
	}

	log.Println("Tables created successfully")
	return nil
}

// InsertPage inserts a new page or increments the view count if the page already exists.
func (d *Database) InsertPage(url string) (int64, error) {
	result, err := d.db.Exec("INSERT INTO pages (url) VALUES (?) ON DUPLICATE KEY UPDATE view_count = view_count + 1", url)
	if err != nil {
		return 0, fmt.Errorf("failed to insert page: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return id, nil
}

// InsertPageView adds a new page view entry.
func (d *Database) InsertPageView(pageID int64, userAgent, referrer, ipAddress string) error {
	_, err := d.db.Exec("INSERT INTO page_views (page_id, user_agent, referrer, ip_address) VALUES (?, ?, ?, ?)", pageID, userAgent, referrer, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to insert page view: %v", err)
	}

	return nil
}
