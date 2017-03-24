package template

import (
	"testing"
	"time"
)

// TestShortDuration tests the shortDuration function
func TestShortDuration(t *testing.T) {
	tests := []struct {
		duration       time.Duration
		expectedString string
	}{
		// Test the empty duration
		{
			duration:       0,
			expectedString: "0s",
		},

		// Test a duration with only seconds
		{
			duration:       5 * time.Second,
			expectedString: "5s",
		},

		// Test a duration with only minutes
		{
			duration:       5 * time.Minute,
			expectedString: "5m",
		},

		// Test a duration with minutes and seconds
		{
			duration:       5*time.Minute + 1*time.Second,
			expectedString: "5m1s",
		},

		// Test a duration with only hours
		{
			duration:       5 * time.Hour,
			expectedString: "5h",
		},

		// Test a duration with hours and minutes
		{
			duration:       5*time.Hour + 5*time.Minute,
			expectedString: "5h5m",
		},
	}

	for _, test := range tests {
		returnedString := shortDuration(test.duration)
		if returnedString != test.expectedString {
			t.Fatalf("expected '%s', returned '%s'", test.expectedString, returnedString)
		}
	}
}
