package fluxgenerator

import (
	"encoding/json"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const (
	fluxGeneratorApiVersion = "generators.giantswarm.io/v1"
	fluxGeneratorKind       = "Konfigure"
	annotationKey           = "config.kubernetes.io/function"
	annotationValue         = `exec:
  path: /plugins/konfigure`
)

type fluxGenerator struct {
	ApiVersion string `json:"api_version,omitempty"`
	Kind       string `json:"kind,omitempty"`

	Metadata fluxGeneratorMetadata `json:"metadata,omitempty"`

	Name                    string `json:"name,omitempty"`
	AppCatalog              string `json:"app_catalog,omitempty"`
	AppDestinationNamespace string `json:"app_destination_namespace,omitempty"`
	AppName                 string `json:"app_name,omitempty"`
	AppVersion              string `json:"app_version,omitempty"`
}

type fluxGeneratorMetadata struct {
	Name        string            `json:"name,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func runFluxGeneratorError(cmd *cobra.Command, args []string) error {
	err := validateFlags()
	if err != nil {
		return microerror.Mask(err)
	}

	generator := fluxGenerator{
		ApiVersion: fluxGeneratorApiVersion,
		Kind:       fluxGeneratorKind,
		Metadata: fluxGeneratorMetadata{
			Name: flag.Name,
			Annotations: map[string]string{
				annotationKey: annotationValue,
			},
		},

		Name:                    flag.Name,
		AppCatalog:              flag.AppCatalog,
		AppDestinationNamespace: flag.AppDestinationNamespace,
		AppName:                 flag.AppName,
		AppVersion:              flag.AppVersion,
	}

	var data []byte

	switch flag.Output {
	case "yaml":
		data, err = yaml.Marshal(generator)
		if err != nil {
			return microerror.Mask(err)
		}
	case "json":
		data, err = json.MarshalIndent(generator, "", "    ")
		if err != nil {
			return microerror.Mask(err)
		}
		data = append(data, '\n')
	default:
		return microerror.Maskf(executionFailedError, "unknown output format %q", flag.Output)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s", data)
	return nil
}
