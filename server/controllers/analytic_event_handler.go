package controllers

import (
	"errors"
	"fmt"
	"github/febzey/go-analytics/types"
	"time"
)

type AnalyticEventHandler struct {

	// Type of event: load | pushsate | onhashchange | popstate | unload | buttonclick
	event string

	// The function that will handle the specified event.
	handler func(AnalyticsPayload, ClientDetails) error
}

func (c *Controller) newEventHandler() {
	events := []AnalyticEventHandler{
		{
			event:   "load",
			handler: c.handleLoadPayload,
		},
		{
			event:   "pushstate",
			handler: c.handleLoadPayload,
		},
		{
			event:   "onhashchange",
			handler: c.handleLoadPayload,
		},
		{
			event:   "unload",
			handler: c.handleUnloadPayload,
		},
		{
			event:   "buttonclick",
			handler: handleButtonClickPayload,
		},
	}

	for _, event := range events {
		c.AnalyticEventHandlers[event.event] = event
	}
}

// Fires for analytic events, for further handling
func (c *Controller) handleAnalyticEvent(payload AnalyticsPayload, clientDetails ClientDetails) error {

	// Getting the handler for specified analytic event
	eventHandler, exists := c.AnalyticEventHandlers[payload.Event]
	if !exists {
		return errors.New("no analytic event found")
	}

	fmt.Println(payload.Event, " analytic event!")

	// Sending our analytic payload to the correct handler function
	err := eventHandler.handler(payload, clientDetails)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

// Fires whenever a client loads or navigates through pages
func (c *Controller) handleLoadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {

	// !! Need a way to not call this first,
	// ! Perhaps use the URL as a unique identifier?
	// ! That way we dont need the PageID before inserting a pageView
	pageID, err := c.db.InsertPage(payload.ClientData.URL)
	if err != nil {
		fmt.Printf("Failed to insert page: %v", err)
		return err
	}

	// we now have a page view struct with a 0 duration time, we add this to cache
	// to be updated later
	pageView := newPageView(pageID, payload, clientDetails)

	if payload.Event != "load" {

		// Update the very last page view's view duration in pageView Cache
		// returns the last page view, along with our page view time
		// ready to be saved to database now.
		lastPageView, err := c.UpdateClientPageViewDuration(payload.ClientData.Token)
		if err != nil {
			return errors.New("error updating page view duration in event handler")
		}

		fmt.Printf("Token: %s, Viewed for %d Seconds", payload.ClientData.Token, lastPageView.ViewDuration)

		//c.savePageAndPageView(payload, clientDetails)

		//! Work on saving page views after the user has left the page
	}

	// Adding the page view to a cache.
	c.AddPageViewToCache(pageView)
	return nil
}

// Analytic event handler for when a client unloads a page.
func (c *Controller) handleUnloadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	_, err := c.UpdateClientPageViewDuration(payload.ClientData.Token)
	if err != nil {
		return err
	}

	// if _, err := c.savePageAndPageView(payload, clientDetails); err != nil {
	// 	return err
	// }
	return nil
}

// Analytic
func handleButtonClickPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	//!Implement me
	return nil
}

// func (c *Controller) savePageAndPageView(payload AnalyticsPayload, clientDetails ClientDetails) (types.PageView, error) {

// 	// lets get that view duration.
// 	pageID, err := c.db.InsertPage(payload.ClientData.URL)
// 	if err != nil {
// 		fmt.Printf("Failed to insert page: %v", err)
// 		return types.PageView{}, err
// 	}

// 	pageView, exists := c.GetLastPageViewByTokenFromCache(payload.ClientData.Token)
// 	if !exists {
// 		return types.PageView{}, errors.New("no page view found")
// 	}

// 	view := types.PageView{
// 		PageID:         pageID,
// 		AnalyticsToken: payload.ClientData.Token,
// 		DeviceWidth:    payload.ClientData.DeviceWidth,
// 		DeviceHeight:   payload.ClientData.DeviceHeight,
// 		UserAgent:      payload.ClientData.UserAgent,
// 		Referrer:       payload.ClientData.Referrer,
// 		Timestamp:      time.Now(),
// 		IPAddress:      clientDetails.IP,
// 		Hostname:       clientDetails.Hostname,
// 		City:           clientDetails.City,
// 		Region:         clientDetails.Region,
// 		Country:        clientDetails.Country,
// 		Loc:            clientDetails.Loc,
// 		Org:            clientDetails.Org,
// 		Postal:         clientDetails.Postal,
// 		Timezone:       clientDetails.Timezone,
// 		ViewDuration:   pageView.ViewDuration,
// 	}

// 	if err := c.db.InsertPageView(view); err != nil {
// 		return types.PageView{}, err
// 	}

// 	return view, nil
// }

// This function will be called at page load, or route change,
// and is immedietly sent to cache to be used when client is done with page.
// we update viewDuration after the user closes or changes pages.
func newPageView(pageID int64, payload AnalyticsPayload, clientDetails ClientDetails) types.PageView {
	view := types.PageView{
		PageID:         pageID,
		URL:            payload.ClientData.URL,
		AnalyticsToken: payload.ClientData.Token,
		DeviceWidth:    payload.ClientData.DeviceWidth,
		DeviceHeight:   payload.ClientData.DeviceHeight,
		UserAgent:      payload.ClientData.UserAgent,
		Referrer:       payload.ClientData.Referrer,
		Timestamp:      time.Now(),
		IPAddress:      clientDetails.IP,
		Hostname:       clientDetails.Hostname,
		City:           clientDetails.City,
		Region:         clientDetails.Region,
		Country:        clientDetails.Country,
		Loc:            clientDetails.Loc,
		Org:            clientDetails.Org,
		Postal:         clientDetails.Postal,
		Timezone:       clientDetails.Timezone,
		ViewDuration:   0,
	}
	return view
}
