package workflow

// GoTest runs `go test`.
type GoTest struct {
}

func (t *GoTest) Run(e Executor) error {
	if err := e.RunCommand("go", "test", "-v", "$(glide novendor)"); err != nil {
		return err
	}

	return nil
}
