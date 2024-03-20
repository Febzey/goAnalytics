package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// client http handlers.

type NewClientRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Adding a new client to our service.
func (c *Controller) PostNewClient(w http.ResponseWriter, r *http.Request) {

	var req NewClientRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		// Handle decoding error
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Println("New client request: ", req)

	// do somethn with the reqW

	err = c.AuthService.CreateAndSaveClient(req.Email, req.Password)
	if err != nil {
		errMsg := fmt.Errorf("error creating client: %v", err)

		if err.Error() == "exists" {
			// return exists status code
			http.Error(w, "Client already exists.", http.StatusConflict)
			return
		}

		fmt.Println(errMsg)
		http.Error(w, "Error creating client.", http.StatusInternalServerError)
		return
	}

	// return to client that it was a success.

	// log some things telling us we created a client
	// and what their username is.
	fmt.Println("Client created successfully: ", req.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// Getting an existing client
// signing in a client for our frontend stuff.
func (c *Controller) GetClient(w http.ResponseWriter, r *http.Request) {
	var req NewClientRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Println("Get client request: ", req)

	client, err := c.AuthService.GetAndVerifyClient(req.Email, req.Password)
	if err != nil {
		fmt.Println("Error getting client: ", err)

		if err.Error() == "not found" {
			http.Error(w, "Client not found.", http.StatusNotFound)
			return
		}

		// if we set a different status code for a bad password,
		// people can attack us by enumerating emails.
		// to find emails we have in our database.
		// so we keep it the same as not found.
		if err.Error() == "badpass" {
			http.Error(w, "Bad password.", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log some things telling us we got a client
	// and what their username is.
	fmt.Println("Client found and authed by password: ", client.Email)

	// no need to send the password back over the network!
	// even if its encrypted.
	client.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// * the key is used to verify the hostname of the website that is sending us analytics data.
// * it will be assigned to the client from the moment its made,
// * but only gets assigned to a hostname after the script is loaded for the first time.
// * after it is loaded for first time, its locked to that hostname, and cannot be used by another hostname.

// * should really just use the client key to create a new key, not the user and password.
func (c *Controller) GetNewPublicKey(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	password := r.Header.Get("Password")

	pubKey, err := c.AuthService.CreateAndSavePublicKey(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pubKey)
}
