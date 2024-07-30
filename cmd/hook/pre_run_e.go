package hook

import (
	"context"

	"github.com/pkg/errors"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

func PreRunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var err error

	var repo *gitrepo.Repo
	{
		cwd := cmd.Flag("working-directory").Value.String()
		dir, err := gitrepo.TopLevel(ctx, cwd)
		if err != nil {
			return microerror.Mask(err)
		}

		c := gitrepo.Config{
			Dir: dir,
		}

		repo, err = gitrepo.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	defaultSha, err := repo.HeadSHA(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	defaultTag, err := repo.HeadTag(ctx)
	if errors.Is(err, &gitrepo.ReferenceNotFoundError{}) {
		defaultTag = ""
	} else if err != nil {
		return microerror.Mask(err)
	}

	defaultBranch, err := repo.HeadBranch(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	// Define the version we are building.
	gitVersion, err := repo.ResolveVersion(ctx, "HEAD")
	if err != nil {
		return microerror.Mask(err)
	}

	cmd.PersistentFlags().String("branch", defaultBranch, "git branch being built")
	cmd.PersistentFlags().String("sha", defaultSha, "git SHA1 being built")
	cmd.PersistentFlags().String("tag", defaultTag, "git tag being built")
	cmd.PersistentFlags().String("version", gitVersion, "version found in git")

	return nil
}
