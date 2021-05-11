package argoapp

import (
	"encoding/json"
	"fmt"

	"github.com/giantswarm/argoapp/pkg/argoapp"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func runAppCRError(cmd *cobra.Command, args []string) error {
	var err error

	config := argoapp.ApplicationConfig{
		Name: flag.Name,

		AppName:                 flag.AppName,
		AppVersion:              flag.AppVersion,
		AppCatalog:              flag.AppCatalog,
		AppDestinationNamespace: flag.AppDestinationNamespace,

		ConfigRef:           flag.ConfigRef,
		DisableForceUpgrade: flag.DisableForceUpgrade,
	}

	applicationCR, err := argoapp.NewApplication(config)
	if err != nil {
		return microerror.Mask(err)
	}

	var data []byte

	switch flag.Output {
	case "yaml":
		data, err = yaml.Marshal(applicationCR)
		if err != nil {
			return microerror.Mask(err)
		}
	case "json":
		data, err = json.MarshalIndent(applicationCR, "", "    ")
		if err != nil {
			return microerror.Mask(err)
		}
		data = append(data, '\n')
	default:
		return microerror.Maskf(executionFailedError, "unknown output format %q", flag.Output)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s", data)
	return nil
}
