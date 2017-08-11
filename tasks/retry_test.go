package tasks

import (
	"fmt"
	"testing"

	"github.com/cenk/backoff"
)

func Test_Task_Retry(t *testing.T) {
	tests := []struct {
		retryTask       RetryTask
		expectedRetries int
	}{
		{
			retryTask: NewRetryTask(
				&backoff.ZeroBackOff{},
				&errorCountTask{
					errorCount: 0,
				},
			),
			expectedRetries: 0,
		},
		{
			retryTask: NewRetryTask(
				&backoff.ZeroBackOff{},
				&errorCountTask{
					errorCount: 1,
				},
			),
			expectedRetries: 1,
		},
		{
			retryTask: NewRetryTask(
				&backoff.ZeroBackOff{},
				&errorCountTask{
					errorCount: 4,
				},
			),
			expectedRetries: 4,
		},
	}

	for index, test := range tests {
		err := test.retryTask.Run()
		if err != nil {
			t.Fatal("test", index+1, "expected", nil, "got", err)
		}

		errorCountTask := test.retryTask.Task.(*errorCountTask)
		if errorCountTask.retries != test.expectedRetries {
			t.Fatal("test", index+1, "expected", test.expectedRetries, "got", errorCountTask.retries)
		}
	}
}

type errorCountTask struct {
	errorCount int
	retries    int
}

func (t *errorCountTask) Run() error {
	if t.errorCount > t.retries {
		t.retries++
		return fmt.Errorf("test error")
	}

	return nil
}

func (t *errorCountTask) Name() string {
	return "errorCount"
}

func (t *errorCountTask) String() string {
	return ""
}
