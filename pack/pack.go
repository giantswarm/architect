package pack

import (
	"fmt"
	"path/filepath"

	"k8s.io/helm/pkg/chartutil"

	"github.com/giantswarm/microerror"
)

const (
	PackageHelmChartTaskName = "package-helm-chart"

	PackageHelmChartTaskString = "%s: %s"
)

type PackageHelmChartTask struct {
	chartDir string
	dst      string
}

// Run package the helm chart at p.chartDir into p.dst.
func (p PackageHelmChartTask) Run() error {
	path, err := filepath.Abs(p.chartDir)
	if err != nil {
		return err
	}

	ch, err := chartutil.LoadDir(path)
	if err != nil {
		return err
	}

	_, err = chartutil.Save(ch, p.dst)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (p PackageHelmChartTask) Name() string {
	return PackageHelmChartTaskName
}

func (p PackageHelmChartTask) String() string {
	return fmt.Sprintf(PackageHelmChartTaskString, p.Name(), p.chartDir)
}

func NewPackageHelmChartTask(chartDir, dst string) PackageHelmChartTask {
	return PackageHelmChartTask{
		chartDir: chartDir,
		dst:      dst,
	}
}
