package create

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/create/appcr"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "create a resource.",
	}
)

func init() {
	Cmd.AddCommand(appcr.Cmd)
}
