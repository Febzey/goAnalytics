package database

import (
	"fmt"
	"github/febzey/go-analytics/utils"
)

// InsertPage inserts a new page or increments the view count if the page already exists.
func (d *Database) InsertPage(url string) (int64, error) {

	mainUrl, route, isSecure, err := utils.ParseURL(url)
	if err != nil {
		return 0, fmt.Errorf("failed to parse url: %v", err)
	}

	var secureEnum int

	if isSecure {
		secureEnum = 1
	} else {
		secureEnum = 0
	}

	result, err := d.db.Exec(`
	INSERT INTO pages (url, route, is_secure) 
	VALUES (?, ?, ?) 
	ON DUPLICATE KEY UPDATE view_count = view_count + 1
	`,
		mainUrl, route, secureEnum,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert page: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return id, nil
}
