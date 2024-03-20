package controllers

import (
	"errors"
	"fmt"
	"github/febzey/go-analytics/types"
	"time"
)

/*
*
* Page View Cache and helper functions.
* Keeping track of page views and updating view times.
*
 */

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

// Get the past page viewed by token.
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

// Updating a view duration for a specific token. Called when user changes the page.
func (c *Controller) UpdateClientPageViewDurationInCache(token string) (types.PageView, error) {
	now := time.Now().Unix()

	lastPageView, exists := c.GetLastPageViewByTokenFromCache(token)
	if !exists {
		return types.PageView{}, errors.New("no token found while updating view duration")
	}

	viewDuration := now - lastPageView.Timestamp.Unix()

	lastPageView.ViewDuration = int(viewDuration)

	return *lastPageView, nil
}

// Doing some checks to determine if the visit to this url is a unique visit.
func (c *Controller) isNewVisit(token, url string) int {
	// Check if the page view is in the cache
	views, ok := c.PageViewCache[token]
	if !ok {
		// Page view not found in cache, check database for uniqueness
		unique, err := c.db.CheckDatabaseForUniqueView(token, url)
		if err != nil {
			fmt.Println("Error checking unique view:", err)
			return 0
		}

		return unique // Return the result directly
	}

	// Check if the URL exists in the cache
	for _, view := range views {
		if view.URL == url {
			return 0
		}
	}

	// URL not found in cache, perform database check for uniqueness
	unique, err := c.db.CheckDatabaseForUniqueView(token, url)
	if err != nil {
		fmt.Println("Error checking unique view:", err)
		return 0
	}

	return unique
}
