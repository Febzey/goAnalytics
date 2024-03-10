package controllers

import (
	"errors"
	"log"
	"time"
)

type AnalyticEventHandler struct {
	event   string
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
			event:   "popstate",
			handler: c.handleLoadPayload,
		},
		{
			event:   "unload",
			handler: handleUnloadPayload,
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
func (c *Controller) handleAnalyticEvent(event string, payload AnalyticsPayload, clientDetails ClientDetails) error {

	eventHandler, exists := c.AnalyticEventHandlers[event]
	if !exists {
		return errors.New("no analytic event found")
	}

	err := eventHandler.handler(payload, clientDetails)
	if err != nil {
		return err
	}

	return nil
}

// Fires whenever a client loads or navigates through pages
func (c *Controller) handleLoadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {

	lastPageView, exists := c.GetLastPageViewByURL(payload.ClientToken, payload.URL)
	if exists {
		timeThreshold := time.Now().Add(-25 * time.Minute).Unix()
		if lastPageView.Time >= timeThreshold {
			return nil
		}
	}

	c.AddPageView(payload.ClientToken, payload.URL)

	pageID, err := c.db.InsertPage(payload.URL)
	if err != nil {
		log.Printf("Failed to insert page: %v", err)
		return err
	}

	pageView := c.createPageView(payload.ClientToken, payload, clientDetails.IP, clientDetails, pageID)

	if err := c.db.InsertPageView(pageView); err != nil {
		log.Printf("Failed to insert page view: %v", err)
	}

	return nil
}

// Analytic event handler for when a client unloads a page.
func handleUnloadPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	//!Implemenent me
	return nil
}

// Analytic
func handleButtonClickPayload(payload AnalyticsPayload, clientDetails ClientDetails) error {
	//!Implement me
	return nil
}
