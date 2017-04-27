package template

import (
	"net/url"
	"testing"
	"time"
)

// Test_listToString tests the listToString filter.
func Test_listToString(t *testing.T) {
	tests := []struct {
		list     []string
		expected string
	}{
		{
			list:     []string{},
			expected: "",
		},
		{
			list: []string{
				"1",
				"2",
				"3",
			},
			expected: "1,2,3",
		},
		{
			list: []string{
				"a",
				"b",
				"c",
			},
			expected: "a,b,c",
		},
		{
			list: []string{
				"foo",
				"b ar",
				"baz",
			},
			expected: "foo,b ar,baz",
		},
		{
			list: []string{
				"",
				" ",
				"",
			},
			expected: ", ,",
		},
	}

	for _, test := range tests {
		e := listToString(test.list)
		if e != test.expected {
			t.Fatalf("expected '%s', got '%s'", test.expected, e)
		}
	}
}

// TestShortDuration tests the shortDuration function.
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

// TestURLString tests the urlString function.
func TestURLString(t *testing.T) {
	tests := []struct {
		URL            url.URL
		expectedString string
	}{
		// Test the empty URL
		{
			URL:            url.URL{},
			expectedString: "",
		},

		// Test a http URL
		{
			URL: url.URL{
				Scheme: "http",
				Host:   "giantswarm.io",
			},
			expectedString: "http://giantswarm.io",
		},

		// Test a https URL
		{
			URL: url.URL{
				Scheme: "https",
				Host:   "giantswarm.io",
			},
			expectedString: "https://giantswarm.io",
		},
	}

	for _, test := range tests {
		returnedString := urlString(test.URL)
		if returnedString != test.expectedString {
			t.Fatalf("expected '%s', returned '%s'", test.expectedString, returnedString)
		}
	}
}
