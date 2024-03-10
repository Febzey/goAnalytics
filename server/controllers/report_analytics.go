package controllers

import (
	"encoding/json"
	"fmt"
	"github/febzey/go-analytics/types"
	"log"
	"net/http"
	"time"
)

//operose

// A structure for incoming analytic data payloads
type AnalyticsPayload struct {

	// The event that triggered the payload, example: "pushstate", "load", "hashchange"
	Event string `json:"event"`

	// The user agent for the client or device
	UserAgent string `json:"userAgent"`

	// The url / page the client triggered the payload from
	URL string `json:"url"`

	// if the payload was a load, there might have been a referrer
	Referrer string `json:"referrer"`

	//! todo - add data for button
	//! Perhaps use this structure for each event, but add a dynamic data struct within this.
	//! like a arbitrary interface struct

	//add api key here eventually.
	//do checks to see if api key owner is for the website its sent from
	//do some permission checks, ex: views counter is free, button clicks is paid or something.

	ClientToken string
}

// * Main controller for handling incoming analytics data,
// * we get the data in the form of URL queries, and return a small GIF image.
// * Updating page views and inserting new page. adding various things to caches like view cache, client details cache,
func (c *Controller) analyticsReportHandler(w http.ResponseWriter, r *http.Request) {

	// URL Queries
	analyticsData := r.URL.Query()

	// Getting Raw analytics data from url query
	// analytic_type := analyticsData["event"]
	// fmt.Println(analytic_type)

	data := analyticsData["data"][0]

	fmt.Println(data)
	// Client token, stored in a clients cookie
	var (
		// Payload for an incoming analytic event, pushstate, load etc.
		payload AnalyticsPayload

		// Token that is retrieved from clients cookie,
		// or a new one is created and stored in a cookie if doesnt exist.
		token string

		// Client details mostly information from IP address.
		// This is for a client that is viewing a page.
		clientDetails ClientDetails
	)

	// populating the payload struct with data from query
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Printf("Failed to unmarshal analytics payload: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Getting IP address for client.
	ip := getClientIP(r)
	// Getting token stored in clients browser cookies.
	token = getAnalyticsToken(r)

	// If there is no token, we create a new token and store it in the clients browser.
	// We also create new client details and store it in cache
	if token == "" {
		token = generateAnalyticsToken()
		setAnalyticsToken(w, token)

		details, err := getNewClientDetails(ip)
		if err != nil {
			fmt.Println("Error: could not get ip details")
		}

		c.updateClientDetails(token, *details)

		clientDetails = *details

	} else {
		// Token exists, so getting it from cache
		details, exists := c.getClientDetails(token)
		if !exists {

			//couldnt find client details in cache,
			//so we make new details and store it to cache.
			newDetails, err := getNewClientDetails(ip)
			if err != nil {
				fmt.Println("Error: could not get ip details")
			}

			c.updateClientDetails(token, *newDetails)

			clientDetails = *newDetails

		} else {

			clientDetails = details

		}

	}

	//should have payload by now idk
	payload.ClientToken = token

	err := c.handleAnalyticEvent(payload.Event, payload, clientDetails)
	if err != nil {
		fmt.Println(" error handling analytic event... ")
	}

	serveFile(w, r)
}

func (c *Controller) createPageView(token string, payload AnalyticsPayload, ip string, clientDetails ClientDetails, pageID int64) types.PageView {
	return types.PageView{
		PageID:         pageID,
		AnalyticsToken: token,
		UserAgent:      payload.UserAgent,
		Referrer:       payload.Referrer,
		Timestamp:      time.Now(),
		IPAddress:      ip,
		IP:             clientDetails.IP,
		Hostname:       clientDetails.Hostname,
		City:           clientDetails.City,
		Region:         clientDetails.Region,
		Country:        clientDetails.Country,
		Loc:            clientDetails.Loc,
		Org:            clientDetails.Org,
		Postal:         clientDetails.Postal,
		Timezone:       clientDetails.Timezone,
	}
}
