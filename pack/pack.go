package pack

import (
	"fmt"
	"log"
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
//
// If p.dst is /foo, and the chart is named bar, with version 1.0.0, this
// will generate /foo/bar-1.0.0.tgz.
func (p PackageHelmChartTask) Run() error {
	path, err := filepath.Abs(p.chartDir)
	if err != nil {
		return microerror.Mask(err)
	}

	// Load chart from a directory.
	ch, err := chartutil.LoadDir(path)
	if err != nil {
		return microerror.Mask(err)
	}

	// Save the chart as an archive in the given directory.
	savedFile, err := chartutil.Save(ch, p.dst)
	if err != nil {
		return microerror.Mask(err)
	}

	log.Printf("chart %s packaged at %s", ch.Metadata.GetName(), savedFile)

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
