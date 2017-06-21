package template

import (
	"encoding/json"
	"net"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/giantswarm/architect/configuration/provider/aws/ec2/instance"
)

var (
	// filters defines functions that can be used in the templates.
	filters = template.FuncMap{
		"ec2InstanceListToString": instance.ListToString,
		"jsonMarshal":             jsonMarshal,
		"ipsToString":             ipsToString,
		"listToString":            listToString,
		"shortDuration":           shortDuration,
		"urlString":               urlString,
	}
)

// listToString takes a string slice and returns a string containing all of its
// items being joined together using a comma.
func listToString(list []string) string {
	return strings.Join(list, ",")
}

// ipsToString takes a net.IP slice and returns a string containing all of its
// items being joined together using a comma.
func ipsToString(IPs []net.IP) string {
	var list []string

	for _, ip := range IPs {
		list = append(list, ip.String())
	}

	return strings.Join(list, ",")
}

// jsonMarshal takes some arbitrary value and applies json.Marshal on it.
func jsonMarshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
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
