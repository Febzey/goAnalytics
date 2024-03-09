package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// A function to disect a URL to obtain the route, and mainUrl as well as secure boolean
func ParseURL(inputUrl string) (host, path string, isSecure bool, err error) {
	// Remove '#' from the input URL
	inputUrl = strings.Replace(inputUrl, "#", "", -1)
	// Use regex to replace consecutive slashes with a single slash
	re := regexp.MustCompile(`([^:/])(/{2,})`)
	inputUrl = re.ReplaceAllString(inputUrl, "$1/")

	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", "", false, err
	}

	// If the scheme is empty, assume http:// as a default
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// If the Host is empty, assume the inputUrl is a relative path
	if u.Host == "" {
		u, err = url.Parse("http://" + inputUrl)
		if err != nil {
			return "", "", false, err
		}
	}

	host = u.Hostname()
	path = u.Path

	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	isSecure = u.Scheme == "https"

	return host, path, isSecure, nil
}
