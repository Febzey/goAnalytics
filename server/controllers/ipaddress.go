package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Getting the clients ip
func getClientIP(r *http.Request) string {
	//Check the X-Forwarded-For header first
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ", ")
		return ips[0]
		//return "38.62.73.45"

	}

	return strings.Split(r.RemoteAddr, ":")[0]
	//return "38.62.73.45"
}

// Getting details on an IP Address
func getNewClientDetails(ipAddress string) (*ClientDetails, error) {

	apiURL := fmt.Sprintf("https://ipinfo.io/%s?token=%s", ipAddress, os.Getenv("ipinfotoken"))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var ClientDetails ClientDetails
	if err := json.NewDecoder(resp.Body).Decode(&ClientDetails); err != nil {
		return nil, err
	}

	return &ClientDetails, nil

}
