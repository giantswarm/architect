package instance

import (
	"testing"
)

// Test_listToString tests the listToString filter.
func Test_listToString(t *testing.T) {
	tests := []struct {
		list     []kind
		expected string
	}{
		{
			list:     []kind{},
			expected: "",
		},
		{
			list: []kind{
				TypeM1Small,
				TypeM1Medium,
				TypeM1Large,
			},
			expected: "m1.small,m1.medium,m1.large",
		},
		{
			list: []kind{
				TypeC1Medium,
				TypeC1XLarge,
			},
			expected: "c1.medium,c1.xlarge",
		},
	}

	for _, test := range tests {
		e := ListToString(test.list)
		if e != test.expected {
			t.Fatalf("expected '%s', got '%s'", test.expected, e)
		}
	}
}
