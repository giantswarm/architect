package hook

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	versionFile = "pkg/project/project.go"
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
	var currentVersion string
	{
		if _, err := os.Stat(versionFile); err == nil {
			filecontents, err := ioutil.ReadFile(versionFile)
			if err != nil {
				return microerror.Mask(err)
			}
			currentVersion, err = getVersionInFile(filecontents)
			if err != nil {
				return microerror.Mask(err)
			}
			currentVersion = fmt.Sprintf("%s-%s", currentVersion, defaultSha)
		} else {
			currentVersion, err = repo.ResolveVersion(ctx, "HEAD")
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	cmd.PersistentFlags().String("branch", defaultBranch, "git branch being built")
	cmd.PersistentFlags().String("sha", defaultSha, "git SHA1 being built")
	cmd.PersistentFlags().String("tag", defaultTag, "git tag being built")
	cmd.PersistentFlags().String("version", currentVersion, "version found in project.go or git")

	return nil
}

func getVersionInFile(filecontents []byte) (string, error) {
	versionRegex := regexp.MustCompile(`(version\s*=\s*)"([0-9]+\.[0-9]+\.[0-9]+)(-dev)?"`)
	currentVersion := versionRegex.FindSubmatch(filecontents)
	if len(currentVersion) < 1 {
		return "", microerror.Maskf(executionFailedError, "file %#q exists but no version was found in it", versionFile)
	}

	return string(currentVersion[2]), nil
}
