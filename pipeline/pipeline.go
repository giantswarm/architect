package pipeline

import (
	"path/filepath"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

func StartChannel(fs afero.Fs, workingDirectory, project string) (string, error) {
	c := &struct {
		Version string `yaml:"version"`
	}{}

	chartFile := filepath.Join(workingDirectory, "helm", project+"-chart", "Chart.yaml")

	y, err := afero.ReadFile(fs, chartFile)
	if err != nil {
		return "", microerror.Mask(err)
	}

	err = yaml.Unmarshal(y, c)
	if err != nil {
		return "", microerror.Mask(err)
	}

	elements := strings.Split(c.Version, ".")
	if len(elements) < 2 {
		return "", microerror.Mask(incorrectChartVersionError)
	}
	// remove initial v
	elements[0] = strings.TrimLeft(elements[0], "v")

	items := append(elements[0:2], initialStability)

	startChannel := strings.Join(items, "-")

	return startChannel, nil
}
