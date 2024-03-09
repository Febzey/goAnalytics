package controllers

import "net/http"

// serve a file
func serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/assets/byone.gif")
}
