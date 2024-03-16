package database

import (
	"fmt"
	"log"
)

func (d *Database) createTables() error {
	// Execute SQL statements to create tables
	statements := []string{
		`CREATE TABLE IF NOT EXISTS pages (
			id INT AUTO_INCREMENT PRIMARY KEY,
			url VARCHAR(255) NOT NULL,
			route VARCHAR(255) NOT NULL,
			is_secure TINYINT(1) DEFAULT 0,
			view_count INT DEFAULT 0,
			unique_view_count INT default 0,
			UNIQUE KEY unique_page (url, route)
		)`,
		`CREATE TABLE IF NOT EXISTS page_views (
			id INT AUTO_INCREMENT PRIMARY KEY,
			page_id INT,
			analytics_token VARCHAR(255),
			device_width INT,
			device_height INT,
			user_agent VARCHAR(255),
			referrer VARCHAR(255),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			ip_address VARCHAR(45),
			city VARCHAR(255),
			region VARCHAR(255),
			country VARCHAR(255),
			loc VARCHAR(255),
			org VARCHAR(255),
			postal VARCHAR(255),
			timezone VARCHAR(255),
			view_duration INT,
			unique_view TINYINT(1) DEFAULT 0,
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
