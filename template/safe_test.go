package template

import "testing"

// TestSafeTemplate tests the SafeTemplate function.
func TestSafeTemplate(t *testing.T) {
	tests := []struct {
		s                    interface{}
		template             string
		expectedContent      string
		expectedErrorMatcher func(error) bool
	}{
		// Test that a nil struct and empty string is templated correctly.
		{
			s:                    nil,
			template:             "",
			expectedContent:      "",
			expectedErrorMatcher: IsNilTemplateStruct,
		},

		// Test a template with an int field.
		{
			s: struct {
				foo int
			}{
				foo: 10,
			},
			template:             "{{ .foo }}",
			expectedContent:      "",
			expectedErrorMatcher: IsNotStringType,
		},

		// Test an empty struct and string is templated correctly.
		{
			s:                    struct{}{},
			template:             "",
			expectedContent:      "",
			expectedErrorMatcher: nil,
		},

		// Test a struct with private fields is templated correctly.
		{
			s: struct {
				private string
			}{
				private: "you-cant-see-me",
			},
			template:             "{{ .private }}",
			expectedContent:      "you-cant-see-me",
			expectedErrorMatcher: nil,
		},

		// Test a template with multiple fields.
		{
			s: struct {
				sha string
			}{
				sha: "foo",
			},
			template:             "{{ .sha }} {{ .sha }}",
			expectedContent:      "foo foo",
			expectedErrorMatcher: nil,
		},
	}

	for index, test := range tests {
		returnedContent, err := SafeTemplate(test.s, []byte(test.template))

		// An error was returned, we expected one, and it's the right one, continue.
		if err != nil && test.expectedErrorMatcher != nil && test.expectedErrorMatcher(err) {
			continue
		}
		// An error was returned, but we did not expect one, fail.
		if err != nil && test.expectedErrorMatcher == nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}
		// An error was not returned, but we expected one, fail.
		if err == nil && test.expectedErrorMatcher != nil {
			t.Fatalf("%v: did not receive expected error\n", index)
		}

		returnedContentString := string(returnedContent)

		if returnedContentString != test.expectedContent {
			t.Fatalf(
				"%v: returned content did not match expected content:\nreturned: %v\nexpected: %v\n",
				index,
				returnedContentString,
				test.expectedContent,
			)
		}
	}
}
