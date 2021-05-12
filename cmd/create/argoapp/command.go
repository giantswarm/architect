package argoapp

import (
	"github.com/spf13/cobra"
)

var flag struct {
	Name                    string
	AppName                 string
	AppVersion              string
	AppDestinationNamespace string
	AppCatalog              string
	ConfigRef               string
	DisableForceUpgrade     bool
	Output                  string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "argoapp",
		Short: "Create Argo Application CR resource.",
		RunE:  runAppCRError,
	}

	cmd.Flags().StringVar(&flag.AppName, "app-name", "", "App name.")
	must(cmd.MarkFlagRequired("app-name"))
	cmd.Flags().StringVar(&flag.AppDestinationNamespace, "app-destination-namespace", "", "Destination namespace where the app should be installed.")
	must(cmd.MarkFlagRequired("app-destination-namespace"))
	cmd.Flags().StringVar(&flag.AppVersion, "app-version", "", "App version.")
	must(cmd.MarkFlagRequired("app-version"))
	cmd.Flags().StringVar(&flag.AppCatalog, "app-catalog", "", "App catalog name.")
	must(cmd.MarkFlagRequired("app-catalog"))
	cmd.Flags().StringVar(&flag.ConfigRef, "config-ref", "", "Configuration version which is a git ref of giantswarm/config repository. Usually major version tag in \"v1\" format.")
	must(cmd.MarkFlagRequired("config-ref"))
	cmd.Flags().BoolVar(&flag.DisableForceUpgrade, "disable-force-upgrade", false, "Disable helm chart force upgrade.")
	cmd.Flags().StringVar(&flag.Name, "name", "", "Generated Argo Application CR name.")
	cmd.Flags().StringVarP(&flag.Output, "output", "o", "yaml", "Output format. Allowed values: yaml, json.")

	return cmd
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
