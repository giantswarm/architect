package appcr

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/app/v3/pkg/app"
)

func runAppCRError(cmd *cobra.Command, args []string) error {
	var err error

	err = validateConfigVersion(flag.ConfigVersion)
	if err != nil {
		return err
	}

	config := app.Config{
		AppName:             flag.AppName,
		AppNamespace:        flag.AppNamespace,
		AppCatalog:          flag.Catalog,
		AppVersion:          flag.AppVersion,
		ConfigVersion:       flag.ConfigVersion,
		DisableForceUpgrade: flag.DisableForceUpgrade,
		Name:                flag.Name,
		UserConfigMapName:   flag.UserConfigMapName,
		UserSecretName:      flag.UserSecretName,
	}

	appCR := app.NewCR(config)

	err = app.Print(os.Stdout, flag.Output, appCR)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
