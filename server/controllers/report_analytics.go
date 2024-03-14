package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//operose

type PageViewData struct {
	//isFirstLoad bool
}

type UnloadPageData struct {
	//view_duration int64
}

type ClientMetaData struct {
	// The user agent for the client or device
	UserAgent string `json:"userAgent"`

	// The url / page the client triggered the payload from
	URL string `json:"url"`

	// if the payload was a load, there might have been a referrer
	Referrer string `json:"referrer"`

	// Device width
	DeviceWidth int `json:"device_width"`

	// Device height
	DeviceHeight int `json:"device_height"`

	// Client token, stored in a cookie.
	Token string
}

type AnalyticsPayload struct {

	// event types can be "load" | "unload" | "pushstate" | "onhashchange"
	Event string `json:"event"`

	// Some metadata on our client viewing. Ex: url, useragent, referrer
	ClientData ClientMetaData `json:"client_meta_data"`

	// Data needed for the event, ex: button data, page view data.
	EventData interface{} `json:"event_data"`
}

// Main controller for handling incoming analytics data,
// we get the data in the form of URL queries, and return a small GIF image.
// Updating page views and inserting new page. adding various things to caches like view cache, client details cache,
func (c *Controller) analyticsReportHandler(w http.ResponseWriter, r *http.Request) {

	var (
		// Payload for an incoming analytic event, pushstate, load etc.
		payload AnalyticsPayload

		// Client details mostly information from IP address.
		// This is for a client that is viewing a page.
		clientDetails ClientDetails
	)

	// URL Queries
	analyticsData := r.URL.Query()

	// Gettig the main analytics payload from our url query
	data := analyticsData.Get("analytics_payload")

	// populating the payload struct with data from query
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Printf("Failed to unmarshal analytics payload: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Getting IP address for client.
	ip := getClientIP(r)
	// Getting token stored in clients browser cookies.
	payload.ClientData.Token = getAnalyticsToken(r)

	// If there is no token, we create a new token and store it in the clients browser.
	// We also create new client details and store it in cache
	if payload.ClientData.Token == "" {
		//Token does not exist, so create a new one and set it as a cookie and save to cache.
		payload.ClientData.Token = generateAnalyticsToken()
		setAnalyticsToken(w, payload.ClientData.Token)

		details, err := getNewClientDetails(ip)
		if err != nil {
			fmt.Println("Error: could not get ip details")
		}

		c.updateClientDetails(payload.ClientData.Token, *details)

		clientDetails = *details

	} else {
		// Token exists, so getting it from cache
		details, exists := c.getClientDetails(payload.ClientData.Token)
		if !exists {

			//couldnt find client details in cache,
			//so we make new details and store it to cache.
			newDetails, err := getNewClientDetails(ip)
			if err != nil {
				fmt.Println("Error: could not get ip details")
			}

			c.updateClientDetails(payload.ClientData.Token, *newDetails)

			clientDetails = *newDetails

		} else {

			clientDetails = details

		}

	}

	//should have payload by now idk

	err := c.handleAnalyticEvent(payload, clientDetails)
	if err != nil {
		fmt.Println("error handling analytic event... ", err)
	}

	serveFile(w, r)
}
