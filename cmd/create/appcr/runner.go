package appcr

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/app/v3/pkg/app"
)

func runAppCRError(cmd *cobra.Command, args []string) error {
	config := app.Config{
		AppName:             flag.AppName,
		AppNamespace:        flag.AppNamespace,
		AppCatalog:          flag.Catalog,
		AppVersion:          flag.AppVersion,
		DisableForceUpgrade: flag.DisableForceUpgrade,
		Name:                flag.Name,
		UserConfigMapName:   flag.UserConfigMapName,
		UserSecretName:      flag.UserSecretName,
	}

	appCR := app.NewCR(config)

	err := app.Print(os.Stdout, flag.Output, appCR)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
