package workflow

import (
	"github.com/giantswarm/architect/pack"
	"github.com/giantswarm/architect/tasks"
	"github.com/spf13/afero"
)

// NewPackageHelmChartTaskCreator returns a closure whith pre-defined dst directory.
//
// This is usefull when used with processHelmDir.
func NewPackageHelmChartTaskCreator(dst string) func(afero.Fs, string, ProjectInfo) (tasks.Task, error) {
	return func(fs afero.Fs, chartDir string, projectInfo ProjectInfo) (tasks.Task, error) {
		return pack.NewPackageHelmChartTask(chartDir, dst), nil
	}
}
