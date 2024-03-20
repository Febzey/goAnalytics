package database

import (
	"fmt"
	"github/febzey/go-analytics/types"
)

/*
*
* Handling individual client page views inside the page_views table.
*
 */

// Inserting a new page view for specific client token.
// Each view is saves as a row seperatly in the page_views table.
func (d *Database) InsertPageView(args types.PageView) error {
	_, err := d.Pool.Exec(`
		INSERT INTO page_views (
			page_id,
			analytics_token,
			device_width,
			device_height,
			user_agent,
			referrer,
			ip_address,
			city,
			region,
			country,
			loc,
			org,
			postal,
			timezone,
			view_duration,
			unique_view,
			url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		args.PageID,
		args.AnalyticsToken,
		args.DeviceWidth,
		args.DeviceHeight,
		args.UserAgent,
		args.Referrer,
		args.IPAddress,
		args.City,
		args.Region,
		args.Country,
		args.Loc,
		args.Org,
		args.Postal,
		args.Timezone,
		args.ViewDuration,
		args.UniqueView,
		args.URL,
	)
	if err != nil {
		return fmt.Errorf("failed to insert page view: %v", err)
	}

	return nil
}

// Checking if a page already exists and has been viewed by a given token and url.
func (d *Database) CheckDatabaseForUniqueView(token, url string) (int, error) {
	var count int

	query := `
		SELECT COUNT(*) FROM page_views WHERE analytics_token = ? AND url = ?
	`

	row := d.Pool.QueryRow(query, token, url)

	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve unique view status: %v", err)
	}

	if count > 0 {
		return 0, nil
	}

	return 1, nil
}
