package helm

import (
	"github.com/giantswarm/architect/cmd/helm/template"
	"github.com/spf13/cobra"
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
