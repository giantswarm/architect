package workflow

type Workflow struct {
	Tasks []Task
}

func NewGolangWorkflow() Workflow {
	workflow := Workflow{
		Tasks: []Task{
			&GoTest{},
			&GoBuild{},
		},
	}

	return workflow
}

func RunWorkflow(executor Executor, workflow Workflow) error {
	for _, task := range workflow.Tasks {
		if err := task.Run(executor); err != nil {
			return err
		}
	}

	return nil
}
