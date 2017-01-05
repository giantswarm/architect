package workflow

// GoBuild runs `go build`.
type GoBuild struct {
}

func (t *GoBuild) Run(e Executor) error {
	if err := e.RunCommand("go", "build", "-a", "-v", "-tags", "netgo"); err != nil {
		return err
	}

	return nil
}
