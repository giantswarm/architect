package appcr

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/app/pkg/app"
)

func runAppCRError(cmd *cobra.Command, args []string) error {
	name := cmd.Flag("name").Value.String()
	appName := cmd.Flag("app-name").Value.String()
	appNamespace := cmd.Flag("app-namespace").Value.String()
	appVersion := cmd.Flag("app-version").Value.String()
	catalog := cmd.Flag("catalog").Value.String()

	disableForceUpgrade, err := cmd.Flags().GetBool("disable-force-upgrade")
	if err != nil {
		return microerror.Mask(err)
	}

	format := cmd.Flag("output").Value.String()

	config := app.Config{
		AppName:             appName,
		AppNamespace:        appNamespace,
		AppCatalog:          catalog,
		AppVersion:          appVersion,
		DisableForceUpgrade: disableForceUpgrade,
		Name:                name,
	}

	appCR := app.NewCR(config)

	err = app.Print(os.Stdout, format, appCR)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
