package main

import (
	"fmt"
	"github/febzey/go-analytics/controllers"
	"github/febzey/go-analytics/internal/database"
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

	db := database.NewDatabase()
	err := db.Init()
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))

	contr := controllers.NewController(router, db)
	contr.LoadRoutes()

	// // Define your routes
	router.PathPrefix("/analytics.js").Handler(http.HandlerFunc(serveAnalytics))

	server := ServerConfig(router)

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server is starting on port: %s \n", os.Getenv("SERVER_PORT"))
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

func serveAnalytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeFile(w, r, "../scriptDevelopment/dist/index.js")
}
