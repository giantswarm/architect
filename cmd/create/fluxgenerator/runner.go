package fluxgenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var (
	configVersionRangeRegexp = regexp.MustCompile(`\d+\.x\.x`)
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

	var configRef string
	if flag.ConfigRef != "" {
		configRef = flag.ConfigRef
	} else {
		dir := strings.TrimSuffix(flag.ConfigRefFromChart, "Chart.yaml")
		path := filepath.Join(dir, "Chart.yaml")
		content, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			return microerror.Mask(fmt.Errorf("file %q does not exist", path))
		}

		var chartYaml struct {
			Annotations map[string]string `json:"annotations"`
		}

		err = yaml.Unmarshal(content, &chartYaml)
		if err != nil {
			return microerror.Mask(fmt.Errorf("failed to parse yaml file %q: %s", path, err))
		}

		annotation := "config.giantswarm.io/version"
		if chartYaml.Annotations == nil || chartYaml.Annotations[annotation] == "" {
			// TODO(kopiczko): When all unique apps are migrated
			// uncomment the code below and delete everything else
			// in this if statement.
			//
			//	return microerror.Mask(fmt.Errorf("annotation %q in file %q not found", annotation, path))
			//
			if chartYaml.Annotations == nil {
				chartYaml.Annotations = map[string]string{}
			}
			chartYaml.Annotations[annotation] = "FIXME"
		}

		configRef = chartYaml.Annotations[annotation]

		if configVersionRangeRegexp.MatchString(configRef) {
			major := strings.SplitN(configRef, ".", 2)[0]
			configRef = "v" + major
		}
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
		AppVersion:              configRef,
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

	fmt.Fprintf(cmd.OutOrStdout(), "%s", data)
	return nil
}
