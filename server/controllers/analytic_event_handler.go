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

	// work on unique visits.

	if payload.Event != "load" {
		lastPageView, err := c.UpdateClientPageViewDuration(payload.ClientData.Token)
		if err != nil {
			fmt.Println(err)
		}

		c.savePageAndPageView(lastPageView)
	}

	pageView := c.newPageView(payload, clientDetails)

	// Adding the page view to a cache.
	c.AddPageViewToCache(pageView)

	return nil
}

// Analytic event handler for when a client unloads a page.
func (c *Controller) handleUnloadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	lastView, err := c.UpdateClientPageViewDuration(payload.ClientData.Token)
	if err != nil {
		return err
	}

	if err := c.savePageAndPageView(lastView); err != nil {
		return err
	}
	return nil
}

// Analytic
func handleButtonClickPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	//!Implement me
	return nil
}

// this is where we will save a page
func (c *Controller) savePageAndPageView(pv types.PageView) error {
	pageID, err := c.db.InsertPage(pv.URL, pv.UniqueView)
	if err != nil {
		fmt.Printf("Failed to insert page: %v", err)
		return err
	}

	pv.PageID = pageID

	// now handle the indivudal page view

	if err := c.db.InsertPageView(pv); err != nil {
		return err
	}

	return nil
}

// This function will be called at page load, or route change,
// and is immedietly sent to cache to be used when client is done with page.
// we update viewDuration after the user closes or changes pages.
// THE PAGE ID IS NOT ASSIGNED HERE, WE ONLY MAKE A VIEW FOR CACHE AND ASSIGN AN ID LATER.
// after this function is called, it is stored to cache. and used again once the user switches pages, or unloads.
func (c *Controller) newPageView(payload AnalyticsPayload, clientDetails ClientDetails) types.PageView {
	var (
		isUnique = 0
	)

	// checking if this page view is unique or not.
	// if the page is unique, we will set it to 1 in the view and store it in cache.
	isPageUnique := c.isNewVisit(payload.ClientData.Token, payload.ClientData.URL)
	if isPageUnique {
		isUnique = 1
	}

	view := types.PageView{
		PageID:         0,
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
		UniqueView:     isUnique,
	}
	return view
}
