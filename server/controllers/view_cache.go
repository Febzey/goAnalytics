package controllers

import (
	"errors"
	"github/febzey/go-analytics/types"
	"time"
)

// type PageView struct {
// 	Time         int64  `json:"time"`
// 	URL          string `json:"url"`
// 	ViewDuration int
// }

// type PageViewCache struct {
// 	views []types.PageView
// }

// Get all page views for a token
func (c *Controller) GetPageViewCache(token string) ([]types.PageView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	views, ok := c.PageViewCache[token]
	return views, ok
}

// Adding a page view to view cache
func (c *Controller) AddPageViewToCache(pv types.PageView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the token already exists in the cache
	if views, ok := c.PageViewCache[pv.AnalyticsToken]; ok {
		// Token exists, add a new page view
		c.PageViewCache[pv.AnalyticsToken] = append(views, pv)
	} else {
		// Token doesn't exist, create a new entry
		c.PageViewCache[pv.AnalyticsToken] = []types.PageView{
			pv,
		}
	}
}

// get last view for url.
func (c *Controller) GetLastPageViewByTokenFromCache(token string) (*types.PageView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if views, ok := c.PageViewCache[token]; ok && len(views) > 0 {
		// Find the very last page view for the specified token based on time
		lastView := &views[0]

		for _, view := range views {
			if view.Timestamp.Unix() > lastView.Timestamp.Unix() {
				lastView = &view
			}
		}

		return lastView, true
	}

	return nil, false
}

// Called when the client is done viewing a page.
// Will update their last viewed pages view duration.
// also returns the last pageView with the updated view duration so we can save it.
func (c *Controller) UpdateClientPageViewDuration(token string) (types.PageView, error) {
	now := time.Now().Unix()

	lastPageView, exists := c.GetLastPageViewByTokenFromCache(token)
	if !exists {
		return types.PageView{}, errors.New("no token found while updating view duration")
	}

	viewDuration := now - lastPageView.Timestamp.Unix()

	lastPageView.ViewDuration = int(viewDuration)

	return *lastPageView, nil
}

// Checking if the given url for a client is a unique page view or not.
func (c *Controller) isNewVisit(token, url string) bool {

	// all page views for this token.
	views, ok := c.PageViewCache[token]

	// unique view!
	if !ok {
		// no views where even found, so we can assume its a unqiue view right??
		return true
	}

	// looping through each page view for token.
	for _, view := range views {

		// we found the same URL that already has the unique tag,
		// meaning there is already a unique view stored for this token and page url.
		if view.URL == url && view.UniqueView == 1 {

			return false

		}

	}

	return true
}
