package release

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
	"github.com/giantswarm/microerror"
)

func runReleaseError(cmd *cobra.Command, args []string) error {
	var projectInfo = workflow.ProjectInfo{
		WorkingDirectory: cmd.Flag("working-directory").Value.String(),
		Organisation:     cmd.Flag("organisation").Value.String(),
		Project:          cmd.Flag("project").Value.String(),

		Sha: cmd.Flag("sha").Value.String(),
		Tag: cmd.Flag("tag").Value.String(),

		Registry:       cmd.Flag("registry").Value.String(),
		DockerUsername: cmd.Flag("docker-username").Value.String(),
		DockerPassword: cmd.Flag("docker-password").Value.String(),

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
		ctx := context.Background()
		tc := oauth2.NewClient(ctx, ts)
		githubClient = github.NewClient(tc)
	}

	var releaseDir string
	{
		path := filepath.Join(cmd.Flag("working-directory").Value.String(), cmd.Flag("destination").Value.String())
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(path, os.ModePerm)
				if err != nil {
					return microerror.Mask(err)
				}
			} else {
				if !os.IsExist(err) {
					return microerror.Mask(err)
				}
			}
		}
		releaseDir = path
	}

	{
		fs := afero.NewOsFs()

		workflow, err := workflow.NewRelease(projectInfo, fs, releaseDir, githubClient, push)
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
