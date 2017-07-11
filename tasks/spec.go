package tasks

import (
	"fmt"
	"log"

	microerror "github.com/giantswarm/microkit/error"
)

// Task represents a piece of work to perform.
type Task interface {
	Name() string
	Run() error

	fmt.Stringer
}

// Run executes a slice of Tasks.
func Run(tasks []Task) error {
	for _, task := range tasks {
		log.Printf("running task: %s\n", task)

		if err := task.Run(); err != nil {
			return microerror.MaskAny(err)
		}
	}

	return nil
}
