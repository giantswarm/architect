package tasks

type DockerTaskConfig struct {
	Volumes          []string
	Env              []string
	WorkingDirectory string
	Image            string
	Args             []string
}

func NewDockerTask(name string, config DockerTaskConfig) ExecTask {
	args := []string{
		"docker", "run",
		"--rm",
	}

	for _, volume := range config.Volumes {
		args = append(args, "-v", volume)
	}

	for _, env := range config.Env {
		args = append(args, "-e", env)
	}

	args = append(args, "-w", config.WorkingDirectory)
	args = append(args, config.Image)
	args = append(args, config.Args...)

	newDockerTask := NewExecTask(name, args)

	return newDockerTask
}
