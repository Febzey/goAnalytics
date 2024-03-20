package database

import (
	"fmt"
	"github/febzey/go-analytics/utils"
)

/*
*
* Handling pages in the database.
* Inserting a new page, or updating an existing pages view count, and unique views.
*
 */

// Inserting a new page or updating existing pages view count or unique views
func (d *Database) InsertPage(url string, isUnique int) (int64, error) {
	mainUrl, route, isSecure, err := utils.ParseURL(url)
	if err != nil {
		return 0, fmt.Errorf("failed to parse url: %v", err)
	}

	secureEnum := 0
	if isSecure {
		secureEnum = 1
	}

	result, err := d.Pool.Exec(`
        INSERT INTO pages (url, route, is_secure, unique_view_count, full_url) 
        VALUES (?, ?, ?, ?, ?) 
        ON DUPLICATE KEY UPDATE view_count = view_count + 1,
                                unique_view_count = unique_view_count + VALUES(unique_view_count)
    `, mainUrl, route, secureEnum, isUnique, url)
	if err != nil {
		return 0, fmt.Errorf("failed to insert page: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return id, nil
}
