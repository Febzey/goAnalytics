package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type AnalyticsData struct {
	Event    string      `json:"event"`
	UserData interface{} `json:"userData"`
	// Add more fields as needed
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	//db := database.NewDatabase()

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))

	// Define your routes
	router.HandleFunc("/analytics", analyticsHandler).Methods("GET")
	router.PathPrefix("/analytics.js").Handler(http.HandlerFunc(serveAnalytics))

	server := ServerConfig(router)

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server is starting on port: %s", os.Getenv("SERVER_PORT"))
		if err := server.ListenAndServe(); err != nil {
			fmt.Sprintln("Server error:", err)
		}
	}()

	// Wait for Ctrl+C signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Start the server
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS by allowing all origins
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Retrieve analytics data from query parameters
	analyticsData := r.URL.Query()

	// Process the analyticsData as needed
	fmt.Printf("Received Analytics Data: %v\n", analyticsData)

	// Respond with a 1x1 transparent pixel (or any other small response)
	http.ServeFile(w, r, "/assets/byone.gif")
}

func serveAnalytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeFile(w, r, "../scriptDevelopment/dist/index.js")
}
