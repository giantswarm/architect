package tasks

// NoopTask is a task which does nothing.
type NoopTask struct {
	Task
}

func (t NoopTask) Run() error {
	return nil
}

func (t NoopTask) Name() string {
	return "noop task"
}

func (t NoopTask) String() string {
	return "noop task"
}
