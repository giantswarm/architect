package circleci

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/circleci/triggerjob"
)

var (
	Cmd = &cobra.Command{
		Use:   "circleci",
		Short: "interacts with CircleCI",
	}
)

func init() {
	Cmd.AddCommand(triggerjob.Cmd)
}
