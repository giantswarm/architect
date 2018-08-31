package tasks

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/giantswarm/microerror"
)

// ExecTask is a task backed by os/exec.
type ExecTask struct {
	name string
	Args []string
}

func (t ExecTask) Run() error {
	cmd := exec.Command(t.Args[0], t.Args[1:]...)

	outputFile := path.Join(os.TempDir(), fmt.Sprintf("architect-%s", t.name))

	f, err := os.Create(outputFile)
	if err != nil {
		return microerror.Mask(err)
	}

	stdoutWriter := io.MultiWriter(os.Stdout, f)
	stderrWriter := io.MultiWriter(os.Stderr, f)

	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	return cmd.Run()
}

func (t ExecTask) Name() string {
	return t.name
}

func (t ExecTask) String() string {
	redactionPhrases := []string{
		"password",
	}

	redactedArgs := []string{}
	for _, arg := range t.Args {
		requiresRedaction := false

		for _, redactionPhrase := range redactionPhrases {
			if strings.Contains(arg, redactionPhrase) {
				requiresRedaction = true
			}
		}

		if requiresRedaction {
			parts := strings.Split(arg, "=")
			if len(parts) == 2 {
				parts[1] = "[REDACTED]"
				arg = parts[0] + "=" + parts[1]
			}
		}

		redactedArgs = append(redactedArgs, arg)
	}

	return fmt.Sprintf("%s:\t'%s'", t.Name(), strings.Join(redactedArgs, " "))
}

func NewExecTask(name string, args []string) ExecTask {
	return ExecTask{
		name: name,
		Args: args,
	}
}
