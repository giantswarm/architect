package argoapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/giantswarm/argoapp/pkg/argoapp"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var (
	configVersionRangeRegexp = regexp.MustCompile(`\d+\.x\.x`)
)

func runAppCRError(cmd *cobra.Command, args []string) error {
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
			return microerror.Mask(fmt.Errorf("annotation %q in file %q not found", annotation, path))
		}

		configRef = chartYaml.Annotations[annotation]

		if configVersionRangeRegexp.MatchString(configRef) {
			major := strings.SplitN(configRef, ".", 2)[0]
			configRef = "v" + major
		}
	}

	config := argoapp.ApplicationConfig{
		Name: flag.Name,

		AppName:                 flag.AppName,
		AppVersion:              flag.AppVersion,
		AppCatalog:              flag.AppCatalog,
		AppDestinationNamespace: flag.AppDestinationNamespace,

		ConfigRef:           configRef,
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
