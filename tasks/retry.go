package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
)

// RetryTask is a task to retry wrapped tasks.
type RetryTask struct {
	BackOff backoff.BackOff
	Task    Task
}

func (t RetryTask) Run() error {
	o := func() error {
		err := t.Task.Run()
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		log.Printf("retrying task '%s' due to error (%s)\n", t.Task.Name(), err.Error())
	}

	err := backoff.RetryNotify(o, t.BackOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (t RetryTask) Name() string {
	return fmt.Sprintf("retry task '%s'", t.Task.Name())
}

func (t RetryTask) String() string {
	return fmt.Sprintf("retry task '%s'", t.Task.String())
}

func NewRetryTask(backOff backoff.BackOff, task Task) RetryTask {
	return RetryTask{
		BackOff: backOff,
		Task:    task,
	}
}
