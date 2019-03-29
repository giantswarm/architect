package release

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
	"github.com/giantswarm/microerror"
)

func runReleaseError(cmd *cobra.Command, args []string) error {
	var err error

	ctx := context.Background()

	fs := afero.NewOsFs()

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: cmd.Flag("working-directory").Value.String(),
		Organisation:     cmd.Flag("organisation").Value.String(),
		Project:          cmd.Flag("project").Value.String(),

		Sha: cmd.Flag("sha").Value.String(),
		Tag: cmd.Flag("tag").Value.String(),

		Version: cmd.Flag("version").Value.String(),
	}

	var githubClient *github.Client
	{
		githubToken := cmd.Flag("github-token").Value.String()
		if githubToken == "" {
			return microerror.Mask(missingGithubTokenError)
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		githubClient = github.NewClient(tc)
	}

	var releaseDir string
	{
		releaseDir, err = ioutil.TempDir(os.TempDir(), "architect-release")
		if err != nil {
			return microerror.Mask(err)
		}
		defer os.RemoveAll(releaseDir)
	}

	{
		workflow, err := workflow.NewRelease(projectInfo, fs, releaseDir, githubClient)
		if err != nil {
			return microerror.Mask(err)
		}

		log.Printf("running workflow: %s\n", workflow)

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return microerror.Mask(err)
		}
		if dryRun {
			log.Printf("dry run, not actually running workflow")
			return nil
		}

		if err := tasks.Run(workflow); err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
