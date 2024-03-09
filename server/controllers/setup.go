package controllers

import (
	"github/febzey/go-analytics/internal/database"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Controller struct {
	// a pointer to our mux router
	r *mux.Router

	// pointer to our database
	db *database.Database

	mu sync.Mutex

	// a cache for our clients,
	// this way we may not have to make as many requests to different service such as ipinfo
	ClientCache map[string]ClientDetails

	PageViews PageViewCache
}

type Route struct {

	//The HTTP method to use for the route.
	Method string

	//The pattern to use for the route.
	Pattern string

	//The handler function to use for the route.
	HandlerFunc http.HandlerFunc
}

func NewController(router *mux.Router, database *database.Database) *Controller {

	c := &Controller{
		r:           router,
		db:          database,
		mu:          sync.Mutex{},
		ClientCache: make(map[string]ClientDetails),
		PageViews:   *NewPageViewCache(),
	}

	return c

}

func (c *Controller) LoadRoutes() {

	var routes = []Route{
		{
			Method:      http.MethodGet,
			Pattern:     "/analytics",
			HandlerFunc: c.analyticsHandler,
		},
		{
			Method:      http.MethodGet,
			Pattern:     "/views",
			HandlerFunc: c.getPageViews,
		},
	}

	for _, route := range routes {
		c.r.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}
}

func (c *Controller) updateClientDetails(token string, details ClientDetails) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ClientCache[token] = details
}

// Function to retrieve user details from the map
func (c *Controller) getClientDetails(token string) (ClientDetails, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	details, exists := c.ClientCache[token]
	return details, exists
}
