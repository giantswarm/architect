package create

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/create/appcr"
	"github.com/giantswarm/architect/cmd/create/argoapp"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "create a resource.",
	}
)

func init() {
	Cmd.AddCommand(appcr.NewCommand())
	Cmd.AddCommand(argoapp.NewCommand())
}
