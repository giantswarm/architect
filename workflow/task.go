package workflow

// Task describes a single step in a workflow.
type Task interface {
	// Run performs the task, returning any errors.
	Run(Executor) error
}
