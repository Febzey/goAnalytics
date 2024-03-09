package controllers

import (
	"encoding/json"
	"fmt"
	"github/febzey/go-analytics/types"
	"log"
	"net/http"
	"time"
)

type AnalyticsPayload struct {

	// The event that triggered the payload, example: "pushstate", "load", "hashchange"
	Event string `json:"event"`

	// The user agent for the client or device
	UserAgent string `json:"userAgent"`

	// The url / page the client triggered the payload from
	URL string `json:"url"`

	// if the payload was a load, there might have been a referrer
	Referrer string `json:"referrer"`

	//add api key here eventually.
	//do checks to see if api key owner is for the website its sent from
	//do some permission checks, ex: views counter is free, button clicks is paid or something.
}

type ClientDetails struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

// TODO: way of authenticating websites who use our service.
// basically if our server detects a new site, we should not do anything with it
// until a real person goes to our website and claims this site, but how do we make sure it is not a bad actor claiming the site?
// what is a good way to implement authentiction for something like this

// our service works by giving a site owner a script src to put in index.html,
// that script talks to our server and sends data over.

// Main controller for handling incoming analytics data,
// we get the data in the form of URL queries, and return a small GIF image.
// Updating page views and inserting new page. adding various things to caches like view cache, client details cache,
// TODO: eventually change to a switch case or something for different events when we start to handle button clicks and stuff
func (c *Controller) analyticsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	analyticsData := r.URL.Query()

	var data string

	// Iterate over the query parameters to find the userData parameter
	for key, values := range analyticsData {
		if key == "data" && len(values) > 0 {
			data = values[0]
			break
		}
	}

	var payload AnalyticsPayload

	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Printf("Failed to unmarshal analytics payload: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// // Respond with a 1x1 transparent pixel (or any other small response)

	// initial page load.
	var token string
	var clientDetails ClientDetails

	ip := getClientIP(r)
	token = getAnalyticsToken(r)

	// if no token present in the cookie.
	// we create a new token, and new client details we then
	// store the token in a cookie and send it to user, then store details in cache
	if token == "" {
		fmt.Println(fmt.Errorf("analytics_token does not exist for this clients browser yet... setting one"))
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

	serveFile(w, r)

	// Determining if we should update database,
	// I think we will only update an existing page once every 20 minutes.
	lastPageView, exists := c.GetLastPageViewByURL(token, payload.URL)
	if exists {
		timeThreshold := time.Now().Add(-25 * time.Minute).Unix()
		if lastPageView.Time >= timeThreshold {
			fmt.Println("not 25 min yet")
		}
	}

	c.AddPageView(token, payload.URL)

	pageID, err := c.db.InsertPage(payload.URL)

	if err != nil {
		log.Printf("Failed to insert page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageView := c.createPageView(token, payload, ip, clientDetails, pageID)
	if err := c.db.InsertPageView(pageView); err != nil {
		log.Printf("Failed to insert page view: %v", err)
		// Handle the error as needed
	}

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
