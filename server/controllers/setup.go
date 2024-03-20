package controllers

import (
	"github/febzey/go-analytics/internal/auth"
	"github/febzey/go-analytics/internal/database"
	"github/febzey/go-analytics/types"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Controller struct {
	// a pointer to our mux router
	r *mux.Router

	// pointer to our database
	db *database.Database

	// A mutex to keep our caches in sync.
	mu sync.Mutex

	// a cache for our clients,
	// this way we may not have to make as many requests to different service such as ipinfo
	ClientCache map[string]ClientDetails

	// a cache for page views.
	PageViewCache map[string][]types.PageView

	// handlers for different analaytic payloads.
	AnalyticEventHandlers map[string]AnalyticEventHandler

	AuthService auth.AuthService
}

type Route struct {

	//The HTTP method to use for the route.
	Method string

	//The pattern to use for the route.
	Pattern string

	//The handler function to use for the route.
	HandlerFunc http.HandlerFunc
}

// creating a new controller,
// Basically the struct that wraps our entire service.
func NewController(router *mux.Router, database *database.Database, authService *auth.AuthService) *Controller {

	c := &Controller{
		r:                     router,
		db:                    database,
		mu:                    sync.Mutex{},
		ClientCache:           make(map[string]ClientDetails),
		PageViewCache:         make(map[string][]types.PageView),
		AnalyticEventHandlers: make(map[string]AnalyticEventHandler),
		AuthService:           *authService,
	}

	// new event handler for incoming analytic events.
	c.newEventHandler()

	return c

}

// Loading HTTP request for our server.
func (c *Controller) LoadRoutes() {

	var routes = []Route{

		// This is the route where our script files will report to.
		{
			Method:      http.MethodGet,
			Pattern:     "/analytics",
			HandlerFunc: c.analyticsReportHandler,
		},

		// Getting page views for a specific URL.
		{
			Method:      http.MethodGet,
			Pattern:     "/views",
			HandlerFunc: c.getPageViews,
		},
		{
			Method:      http.MethodPost,
			Pattern:     "/post/register",
			HandlerFunc: c.PostNewClient,
		},
		{
			Method:      http.MethodPost,
			Pattern:     "/post/login",
			HandlerFunc: c.GetClient,
		},
	}

	for _, route := range routes {
		c.r.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}
}
