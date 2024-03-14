package types

import "time"

// Page represents the pages table.
type Page struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	ViewCount int    `json:"view_count"`
}

// PageView represents the page_views table.
type PageView struct {
	ID             int64     `json:"id"`
	PageID         int64     `json:"page_id"`
	URL            string    `json:"url"`
	DeviceWidth    int       `json:"device_width"`
	DeviceHeight   int       `json:"device_height"`
	UserAgent      string    `json:"user_agent"`
	Referrer       string    `json:"referrer"`
	Timestamp      time.Time `json:"timestamp"`
	IPAddress      string    `json:"ip_address"`
	AnalyticsToken string    `json:"analytics_token"`
	Hostname       string    `json:"hostname"`
	City           string    `json:"city"`
	Region         string    `json:"region"`
	Country        string    `json:"country"`
	Loc            string    `json:"loc"`
	Org            string    `json:"org"`
	Postal         string    `json:"postal"`
	Timezone       string    `json:"timezone"`
	ViewDuration   int       `json:"view_duration"`
}
