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

	cmd.Flags().String("catalog", "", "app catalog name")
	cmd.Flags().String("name", "", "cr name")
	cmd.Flags().String("app-name", "", "app name")
	cmd.Flags().StringP("output", "o", "yaml", "output format. allowed: yaml,json")
	cmd.Flags().String("app-version", "", "app version")

	return cmd
}