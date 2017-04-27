package template

import (
	"net/url"
	"strings"
	"text/template"
	"time"
)

var (
	// filters defines functions that can be used in the templates.
	filters = template.FuncMap{
		"listToString":  listToString,
		"shortDuration": shortDuration,
		"urlString":     urlString,
	}
)

// listToString takes a string slice and returns a string containing all of its
// items being joined together using a comma.
func listToString(list []string) string {
	return strings.Join(list, ",")
}

// shortDuration takes a duration, and provides a shorter string version.
// e.g: Instead of 5m0s, 5m
func shortDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

// urlString takes a URL, and provides a format URL.
// Useful in templates where String is not available
func urlString(u url.URL) string {
	return u.String()
}
