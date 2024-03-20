package controllers

import (
	"errors"
	"fmt"
	"github/febzey/go-analytics/types"
	"strings"
	"time"
)

// *
// *
// * Handling Analytic Event Payloads.
// * here we handle and route page views, button clicks, unloads, etc,
// * anything incoming from our analytics script will be handled and processed here.
// *
// *

// The structure for individual analytic event handlers.
// each analytic event will corrospond to its own handler function and event name.
// event names should very closely resemble the actual event names in javascript that we use.
type AnalyticEventHandler struct {

	// Type of event: load | pushsate | onhashchange | popstate | unload | buttonclick
	event string

	// The function that will handle the specified event.
	handler func(AnalyticsPayload, ClientDetails) error
}

// Registering events for individual analytic event handlers such as page loads, button clicks etc.
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

// This function is called for incoming arbitrary analytic event
// Then the event and the data along with the event is routed to its respective event function handler.
func (c *Controller) handleAnalyticEvent(payload AnalyticsPayload, clientDetails ClientDetails) error {

	eventHandler, exists := c.AnalyticEventHandlers[payload.Event]
	if !exists {
		return errors.New("no analytic event found")
	}

	fmt.Println(payload.Event, " analytic event!")

	err := eventHandler.handler(payload, clientDetails)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

// Fires whenever a client loads or navigates through pages
func (c *Controller) handleLoadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	if payload.Event != "load" {
		lastPageView, err := c.UpdateClientPageViewDurationInCache(payload.ClientData.Token)
		if err != nil {
			fmt.Println(err)
		}

		c.savePageAndPageView(lastPageView)
	}

	pageView := c.newPageView(payload, clientDetails)

	c.AddPageViewToCache(pageView)

	return nil
}

// Analytic event handler for when a client unloads a page.
func (c *Controller) handleUnloadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	lastView, err := c.UpdateClientPageViewDurationInCache(payload.ClientData.Token)
	if err != nil {
		return err
	}

	if err := c.savePageAndPageView(lastView); err != nil {
		return err
	}
	return nil
}

// Analytic Button event, button clicks.
func handleButtonClickPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	//!Implement me
	return nil
}

// Saving page and page_view to database.
func (c *Controller) savePageAndPageView(pv types.PageView) error {
	pageID, err := c.db.InsertPage(pv.URL, pv.UniqueView)
	if err != nil {
		fmt.Printf("Failed to insert page: %v", err)
		return err
	}

	pv.PageID = pageID

	if err := c.db.InsertPageView(pv); err != nil {
		return err
	}

	return nil
}

// Generating structure for a page view, this is created and saved to cache.
// We check if the page is unique here as well as remove any trailing "/" from the URl
func (c *Controller) newPageView(payload AnalyticsPayload, clientDetails ClientDetails) types.PageView {
	payload.ClientData.URL = strings.TrimRight(payload.ClientData.URL, "/")

	isPageUnique := c.isNewVisit(payload.ClientData.Token, payload.ClientData.URL)

	view := types.PageView{
		PageID:         0,
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
		UniqueView:     isPageUnique,
		URL:            payload.ClientData.URL,
	}
	return view
}
