package template

import (
	"strings"
	"time"
)

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
