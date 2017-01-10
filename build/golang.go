package build

import (
	"strings"

	"github.com/giantswarm/architect/command"
)

type GolangBuilder struct {
	buildDir string
}

func (b *GolangBuilder) Test() error {
	out, err := command.RunWithOutput([]string{"glide", "novendor"}, b.buildDir)
	if err != nil {
		return err
	}

	trimmedNovendor := strings.TrimSpace(out)
	splitNoVendor := strings.Split(trimmedNovendor, "\n")
	novendor := strings.Join(splitNoVendor, " ")

	if err := command.Run([]string{"go", "test", "-v", novendor}, b.buildDir); err != nil {
		return err
	}

	return nil
}

func (b *GolangBuilder) Build() error {
	if err := command.Run([]string{"go", "build", "-v", "-a", "-tags", "netgo"}, b.buildDir); err != nil {
		return err
	}

	return nil
}
