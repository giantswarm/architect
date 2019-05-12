package circleci

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/circleci/triggerbuild"
)

var (
	Cmd = &cobra.Command{
		Use:   "circleci",
		Short: "interacts with CircleCI",
	}
)

func init() {
	Cmd.AddCommand(triggerbuild.Cmd)
}
