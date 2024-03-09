package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func getAnalyticsToken(r *http.Request) string {
	cookie, err := r.Cookie("analytics_token")
	if err == nil && cookie != nil {
		// Cookie exists, return its value
		return cookie.Value
	}

	return ""
}

func setAnalyticsToken(w http.ResponseWriter, analyticsToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "analytics_token",
		Value:    analyticsToken,
		Path:     "/",
		MaxAge:   31536000, // 1 year in seconds
		HttpOnly: true,
		Secure:   true, // Set to true if your site is served over HTTPS
	})
}

func generateAnalyticsToken() string {
	tokenBytes := make([]byte, 9)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token
}
