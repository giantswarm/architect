package tasks

import (
	"bytes"

	"github.com/giantswarm/microerror"
	"golang.org/x/sync/errgroup"
)

// ConcurrentTask represent a set of tasks that are executed concurrently using
// an errgroup.Group.
// Externally it act as a Task, and as such, can be included on a workflow or
// wrapped on a retry. Keep in mmind in the latter case that, given the nature
// of errgroup.Group, on failure the whole task set will be retried, including
// tasks that potentially have already finished.
type ConcurrentTask struct {
	name  string
	tasks []Task
}

func NewConcurrentTask(name string, tasks ...Task) *ConcurrentTask {
	return &ConcurrentTask{
		name:  name,
		tasks: tasks,
	}
}

func (c *ConcurrentTask) Run() error {
	var g errgroup.Group

	for _, t := range c.tasks {
		g.Go(t.Run)
	}

	err := g.Wait()

	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (c *ConcurrentTask) Name() string {
	return c.String()
}

func (c *ConcurrentTask) String() string {
	var buffer bytes.Buffer

	for _, t := range c.tasks {
		buffer.WriteString(t.Name())
		buffer.WriteString(";")
	}

	return buffer.String()
}
