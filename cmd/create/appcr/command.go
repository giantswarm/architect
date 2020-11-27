package appcr

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appcr",
		Short: "create App CR resource.",
		RunE:  runAppCRError,
	}

	cmd.Flags().String("app-name", "", "app name")
	cmd.Flags().String("app-namespace", "giantswarm", "app namespace")
	cmd.Flags().String("app-version", "", "app version")
	cmd.Flags().String("catalog", "", "app catalog name")
	cmd.Flags().String("configmap-name", "", "user configmap name")
	cmd.Flags().Bool("disable-force-upgrade", false, "disable helm chart force upgrade")
	cmd.Flags().String("name", "", "cr name")
	cmd.Flags().StringP("output", "o", "yaml", "output format. allowed: yaml,json")
	cmd.Flags().String("secret-name", "", "user secret name")

	return cmd
}
