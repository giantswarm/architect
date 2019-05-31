package template

import (
	"log"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/workflow"
)

func runTemplateError(cmd *cobra.Command, args []string) error {
	var (
		chartDir = cmd.Flag("dir").Value.String()
		sha      = cmd.Flag("sha").Value.String()
		tag      = cmd.Flag("tag").Value.String()
		version  = cmd.Flag("version").Value.String()
	)

	var projectInfo = workflow.ProjectInfo{
		Sha:     sha,
		Tag:     tag,
		Version: version,
	}

	fs := afero.NewOsFs()

	log.Printf("templating helm chart\ndir: %s\nsha: %s\ntag: %s\nversion: %s\n", chartDir, sha, tag, version)

	helmChartTemplate, err := workflow.NewTemplateHelmChartTask(fs, chartDir, projectInfo)
	if err != nil {
		return microerror.Mask(err)
	}

	if err := helmChartTemplate.Run(); err != nil {
		return microerror.Mask(err)
	}

	log.Println("templated helm chart")

	return nil
}
