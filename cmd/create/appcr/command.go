package appcr

import (
	"github.com/spf13/cobra"
)

var flag struct {
	AppName              string
	AppNamespace         string
	AppVersion           string
	Catalog              string
	ConfigVersion        string
	DisableForceUpgrade  bool
	Name                 string
	PauseReconcilliation bool
	Output               string
	UserConfigMapName    string
	UserSecretName       string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appcr",
		Short: "create App CR resource.",
		RunE:  runAppCRError,
	}

	cmd.Flags().StringVar(&flag.AppName, "app-name", "", "app name")
	cmd.Flags().StringVar(&flag.AppNamespace, "app-namespace", "giantswarm", "app namespace")
	cmd.Flags().StringVar(&flag.AppVersion, "app-version", "", "app version")
	cmd.Flags().StringVar(&flag.Catalog, "catalog", "", "app catalog name")
	cmd.Flags().StringVar(&flag.ConfigVersion, "config-version", "", "major version of giantswarm/config to use for this App CR")
	cmd.Flags().BoolVar(&flag.DisableForceUpgrade, "disable-force-upgrade", false, "disable helm chart force upgrade")
	cmd.Flags().StringVar(&flag.Name, "name", "", "CR name")
	cmd.Flags().StringVarP(&flag.Output, "output", "o", "yaml", "output format. allowed: yaml,json")
	cmd.Flags().BoolVar(&flag.PauseReconcilliation, "pause", false, "add pause reconcilliation annotation")
	cmd.Flags().StringVar(&flag.UserConfigMapName, "user-configmap-name", "", "user configmap name")
	cmd.Flags().StringVar(&flag.UserSecretName, "user-secret-name", "", "user secret name")

	return cmd
}
