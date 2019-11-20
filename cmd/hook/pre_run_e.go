package hook

import (
	"context"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

func PreRunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var err error

	var repo *gitrepo.Repo
	{

		dir, err := gitrepo.TopLevel(ctx, ".")
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
	if gitrepo.IsReferenceNotFound(err) {
		defaultTag = ""
	} else if err != nil {
		return microerror.Mask(err)
	}

	defaultBranch, err := repo.HeadBranch(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	// Define the version we are building.
	defaultVersion, err := repo.ResolveVersion(ctx, "HEAD")
	if err != nil {
		return microerror.Mask(err)
	}

	cmd.PersistentFlags().String("branch", defaultBranch, "git branch being built")
	cmd.PersistentFlags().String("sha", defaultSha, "git SHA1 being built")
	cmd.PersistentFlags().String("tag", defaultTag, "git tag being built")
	cmd.PersistentFlags().String("version", defaultVersion, "project version being built")

	return nil
}
