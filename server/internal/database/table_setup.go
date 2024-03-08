package database

import (
	"fmt"
	"log"
)

func (d *Database) createTablesIfNotExist() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Table for page analytics
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS page_analytics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			url VARCHAR(255) UNIQUE,
			view_count INT DEFAULT 0
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create page_analytics table: %v", err)
	}

	// Table for individual page views
	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS page_views (
			id INT AUTO_INCREMENT PRIMARY KEY,
			page_url VARCHAR(255),
			ip_address VARCHAR(45),
			last_view TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create page_views table: %v", err)
	}

	log.Println("Database tables created or already exist")
	return nil
}
