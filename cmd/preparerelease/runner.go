package preparerelease

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/preparerelease/internal"
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

	updateChangelog := cmd.Flag("update-changelog").Value.String()
	if updateChangelog != "" && updateChangelog != "true" && updateChangelog != "false" {
		return microerror.Maskf(executionFailedError, "--update-changelog flag must be true or false")
	}

	var m *internal.Modifier
	{
		c := internal.ModifierConfig{
			NewVersion: version,
			Repo:       repo,
			WorkingDir: workingDir,
		}

		m, err = internal.NewModifier(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	if updateChangelog != "false" {
		err = m.AddReleaseToChangelogMd()
		if err != nil {
			return microerror.Mask(err)
		}
		cmd.Printf("File %#q updated.\n", internal.FileChangelogMd)
	}

	err = m.UpdateVersionInProjectGo()
	if internal.IsFileNotFound(err) {
		// Fall trough. Some projects do not have project.go file.
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		cmd.Printf("File %#q updated.\n", internal.FileProjectGo)
	}

	return nil
}
