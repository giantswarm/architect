package fluxgenerator

import (
	"fmt"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

var flag struct {
	Name                    string
	AppName                 string
	AppVersion              string
	AppDestinationNamespace string
	AppCatalog              string
	DisableForceUpgrade     bool
	Output                  string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fluxgenerator",
		Short: "Create Flux Generator",
		RunE:  runFluxGeneratorError,
	}

	cmd.Flags().StringVar(&flag.AppName, "app-name", "", "App name.")
	cmd.Flags().StringVar(&flag.AppDestinationNamespace, "app-destination-namespace", "", "Destination namespace where the app should be installed.")
	cmd.Flags().StringVar(&flag.AppVersion, "app-version", "", "App version.")
	cmd.Flags().StringVar(&flag.AppCatalog, "app-catalog", "", "App catalog name.")
	cmd.Flags().BoolVar(&flag.DisableForceUpgrade, "disable-force-upgrade", false, "Disable helm chart force upgrade.")
	cmd.Flags().StringVar(&flag.Name, "name", "", "Generated Application CR name.")
	cmd.Flags().StringVarP(&flag.Output, "output", "o", "yaml", "Output format. Allowed values: yaml, json.")

	return cmd
}

func validateFlags() error {
	var errors []string

	if flag.Name == "" {
		errors = append(errors, "--name is required")
	}
	if flag.AppName == "" {
		errors = append(errors, "--app-name is required")
	}
	if flag.AppVersion == "" {
		errors = append(errors, "--app-version is required")
	}
	if flag.AppCatalog == "" {
		errors = append(errors, "--app-catalog is required")
	}
	if flag.AppDestinationNamespace == "" {
		errors = append(errors, "--app-destination-namespace is required")
	}
	if len(errors) != 0 {
		return microerror.Mask(fmt.Errorf("invalid flag(s): %s", strings.Join(errors, ", ")))
	}

	return nil
}
