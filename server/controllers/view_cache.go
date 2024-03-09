package controllers

import "time"

type PageView struct {
	Time int64  `json:"time"`
	URL  string `json:"url"`
}
type PageViewCache struct {
	PageViews map[string][]PageView `json:"pageViews"`
}

func NewPageViewCache() *PageViewCache {
	return &PageViewCache{
		PageViews: make(map[string][]PageView),
	}
}

// Get all page views for a token
func (c *Controller) GetPageViews(token string) ([]PageView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	views, ok := c.PageViews.PageViews[token]
	return views, ok
}

// Adding a page view to view cache
func (c *Controller) AddPageView(token, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the token already exists in the cache
	if views, ok := c.PageViews.PageViews[token]; ok {
		// Token exists, add a new page view
		c.PageViews.PageViews[token] = append(views, PageView{
			Time: time.Now().Unix(),
			URL:  url,
		})
	} else {
		// Token doesn't exist, create a new entry
		c.PageViews.PageViews[token] = []PageView{
			{
				Time: time.Now().Unix(),
				URL:  url,
			},
		}
	}
}

// get last view for url.
func (c *Controller) GetLastPageViewByURL(token, url string) (PageView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if views, ok := c.PageViews.PageViews[token]; ok {
		// Find the last page view for the specified URL
		var lastView PageView

		found := false

		for _, view := range views {
			if view.URL == url && (!found || view.Time > lastView.Time) {
				lastView = view
				found = true
			}
		}

		return lastView, found
	}

	return PageView{}, false
}
