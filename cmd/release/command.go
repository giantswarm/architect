package release

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/release/chart"
	"github.com/giantswarm/architect/cmd/release/prepare"
)

var (
	Cmd = &cobra.Command{
		Use:   "release",
		Short: "release operator versions or charts",
	}
)

func init() {
	Cmd.AddCommand(chart.Cmd)
	Cmd.AddCommand(prepare.Cmd)
}
