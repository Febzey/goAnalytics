package database

import (
	"fmt"
	"github/febzey/go-analytics/utils"
)

// InsertPage inserts a new page or increments the view count if the page already exists.
// isUnique 1 == true | 0 == false
func (d *Database) InsertPage(url string, isUnique int) (int64, error) {
	mainUrl, route, isSecure, err := utils.ParseURL(url)
	if err != nil {
		return 0, fmt.Errorf("failed to parse url: %v", err)
	}

	secureEnum := 0
	if isSecure {
		secureEnum = 1
	}

	result, err := d.db.Exec(`
        INSERT INTO pages (url, route, is_secure, unique_view_count) 
        VALUES (?, ?, ?, ?) 
        ON DUPLICATE KEY UPDATE view_count = view_count + 1,
                                unique_view_count = unique_view_count + VALUES(unique_view_count)
    `, mainUrl, route, secureEnum, isUnique)
	if err != nil {
		return 0, fmt.Errorf("failed to insert page: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return id, nil
}
