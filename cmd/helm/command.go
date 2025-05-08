package helm

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/helm/template"
)

var (
	Cmd = &cobra.Command{
		Use:   "helm",
		Short: "manages helm charts",
	}
)

func init() {
	Cmd.AddCommand(template.Cmd)
}
