package project

import (
	"github.com/giantswarm/architect/cmd/project/version"
	"github.com/spf13/cobra"
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
