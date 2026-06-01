package create

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/v2/cmd/create/appcr"
	"github.com/giantswarm/architect/v2/cmd/create/fluxgenerator"
	"github.com/giantswarm/architect/v2/cmd/create/kustomization"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "create a resource.",
	}
)

func init() {
	Cmd.AddCommand(appcr.NewCommand())
	Cmd.AddCommand(fluxgenerator.NewCommand())
	Cmd.AddCommand(kustomization.NewCommand())
}
