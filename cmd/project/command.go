package project

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/v2/cmd/project/version"
)

var (
	Cmd = &cobra.Command{
		Use:   "project",
		Short: "show project informations",
	}
)

func init() {
	Cmd.AddCommand(version.Cmd)
}
