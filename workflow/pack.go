package workflow

import (
	"github.com/spf13/afero"

	"github.com/giantswarm/architect/pack"
	"github.com/giantswarm/architect/tasks"
)

// NewPackageHelmChartTaskCreator returns a closure with pre-defined dst directory.
//
// This is usefull when used with processHelmDir.
func NewPackageHelmChartTaskCreator(dst string) func(afero.Fs, string, ProjectInfo) (tasks.Task, error) {
	return func(fs afero.Fs, chartDir string, projectInfo ProjectInfo) (tasks.Task, error) {
		return pack.NewPackageHelmChartTask(chartDir, dst), nil
	}
}
