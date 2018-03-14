package tasks

import (
	"bytes"

	"github.com/giantswarm/microerror"
	"golang.org/x/sync/errgroup"
)

type ConcurrentTask struct {
	name  string
	tasks []Task
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

func NewConcurrentTask(name string, tasks ...Task) *ConcurrentTask {
	return &ConcurrentTask{
		name:  name,
		tasks: tasks,
	}
}
