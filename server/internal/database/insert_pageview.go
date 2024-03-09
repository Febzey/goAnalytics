package database

import (
	"fmt"
	"github/febzey/go-analytics/types"
)

// InsertPageView adds a new page view entry.
func (d *Database) InsertPageView(args types.PageView) error {
	_, err := d.db.Exec(`
		INSERT INTO page_views (
			page_id,
			analytics_token,
			user_agent,
			referrer,
			ip_address,
			city,
			region,
			country,
			loc,
			org,
			postal,
			timezone
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		args.PageID,
		args.AnalyticsToken,
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
	)
	if err != nil {
		return fmt.Errorf("failed to insert page view: %v", err)
	}

	return nil
}