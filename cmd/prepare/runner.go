package prepare

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	changelogFile = "CHANGELOG.md"
	versionFile   = "pkg/project/project.go"
)

func runPrepareRelease(cmd *cobra.Command, args []string) error {
	var err error

	workingDir := cmd.Flag("working-directory").Value.String()

	var repo string
	{
		o := cmd.Flag("organisation").Value.String()
		p := cmd.Flag("project").Value.String()
		repo = o + "/" + p
	}

	version := cmd.Flag("version").Value.String()
	if version == "" {
		return microerror.Maskf(executionFailedError, "--version flag can't be empty")
	}

	var m *modifier
	{
		c := modifierConfig{
			NewVersion: version,
			Repo:       repo,
			WorkingDir: workingDir,
		}

		m, err = newModifier(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	err = m.AddReleaseToChangelogMd()
	if err != nil {
		return microerror.Mask(err)
	}

	err = m.UpdateVersionInProjectGo()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
