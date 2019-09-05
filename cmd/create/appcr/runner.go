package appcr

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/pkg/app"
)

func runAppCRError(cmd *cobra.Command, args []string) error {
	name := cmd.Flag("name").Value.String()
	appName := cmd.Flag("app-name").Value.String()
	appVersion := cmd.Flag("app-version").Value.String()
	catalog := cmd.Flag("catalog").Value.String()
	format := cmd.Flag("output").Value.String()

	appCR := app.NewCR(name, appName, appVersion, catalog)

	err := app.Print(os.Stdout, format, appCR)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
