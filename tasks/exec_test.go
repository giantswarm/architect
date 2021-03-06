package tasks

import (
	"testing"
)

func TestExecTaskString(t *testing.T) {
	tests := []struct {
		execTask       ExecTask
		expectedString string
	}{
		{
			execTask: NewExecTask(
				"docker-run",
				[]string{"docker", "run"},
			),
			expectedString: "docker-run:\t'docker run'",
		},
		{
			execTask: NewExecTask(
				"docker-login",
				[]string{"docker", "login", "--password=bar"},
			),
			expectedString: "docker-login:\t'docker login --password=[REDACTED]'",
		},
		{
			execTask: NewExecTask(
				"many-pass",
				[]string{"foo", "--first-password=bar", "--second-password=baz"},
			),
			expectedString: "many-pass:\t'foo --first-password=[REDACTED] --second-password=[REDACTED]'",
		},
		{
			execTask: NewExecTask(
				"boolean-flag",
				[]string{"foo", "-password"},
			),
			expectedString: "boolean-flag:\t'foo -password'",
		},
	}

	for index, test := range tests {
		returnedString := test.execTask.String()

		if returnedString != test.expectedString {
			t.Fatalf(
				"%v: expected string did not match returned\nexpected:\n%s\nreturned: \n%s\n",
				index,
				test.expectedString,
				returnedString,
			)
		}
	}
}
