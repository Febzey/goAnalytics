package controllers

import (
	"encoding/json"
	"github/febzey/go-analytics/utils"
	"log"
	"net/http"
)

type Response struct {
	URL       string `json:"url"`
	Route     string `json:"route"`
	ViewCount int    `json:"view_count"`
}

//TODO: handle all possible cases and formats for domains and urls.

// Queries: ?url="localhost/#blogs/"
// example: http://localhost:8080/views?url="localhost/#/blogs/"
func (c *Controller) getPageViews(w http.ResponseWriter, r *http.Request) {

	url, route, _, err := utils.ParseURL(r.URL.Query().Get("url"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	viewCount, err := c.db.SelectPageRouteViews(url, route)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a response struct
	response := Response{
		URL:       url,
		Route:     route,
		ViewCount: viewCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
