package database

import "fmt"

func (d *Database) SelectPageRouteViews(url, route string) (int, error) {

	var viewCount int

	err := d.db.QueryRow("SELECT view_count from pages WHERE url = ? AND route = ?", url, route).Scan(&viewCount)
	if err != nil {
		return 0, fmt.Errorf("failed to select page route views: %v", err)
	}

	return viewCount, nil
}
