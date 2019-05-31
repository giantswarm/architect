package hook

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

func PreRunE(cmd *cobra.Command, args []string) error {
	// Use git HEAD as defaultSha.
	var defaultSha string
	{
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			return microerror.Maskf(gitNoSHAError, "could not get git sha: %#v\n", err)
		}
		defaultSha = strings.TrimSpace(string(out))
	}

	// Use git tag when available.
	var defaultTag string
	{
		out, err := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD").Output()
		if err == nil {
			// Always populate tag unless building on CircleCI and it is not explicitly requesting to build a tag.
			if _, ciTagExists := os.LookupEnv("CIRCLE_TAG"); os.Getenv("CIRCLECI") != "true" || ciTagExists {
				defaultTag = strings.TrimPrefix(strings.TrimSpace(string(out)), "v")
			}
		}
	}

	// We also use the git HEAD branch as well.
	var defaultBranch string
	{
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			return microerror.Maskf(gitNoBranchError, "could not get git branch: %#v\n", err)
		}
		defaultBranch = strings.TrimSpace(string(out))
	}

	// Define the version we are building.
	var defaultVersion string
	{
		// version can be of three different formats:
		//   2.0.0: building a tagged version.
		//   2.0.0-3a955cbb126f0fe5d51aedf2eb84acca7b074374: building ahead of a tagged version.
		//   1.0.0-939f5c6949f83c0a7ea98a25bc9524fd2f751ffe: building a repo which has no tags.
		if defaultTag != "" {
			defaultVersion = defaultTag
		} else {
			out, err := exec.Command("git", "describe", "--tags", "--abbrev=0", "HEAD").Output()
			if err != nil {
				defaultVersion = fmt.Sprintf("1.0.0-%s", defaultSha)
			} else {
				defaultVersion = fmt.Sprintf("%s-%s", strings.TrimPrefix(strings.TrimSpace(string(out)), "v"), defaultSha)
			}

		}
	}

	cmd.PersistentFlags().String("branch", defaultBranch, "git branch being built")
	cmd.PersistentFlags().String("sha", defaultSha, "git SHA1 being built")
	cmd.PersistentFlags().String("tag", defaultTag, "git tag being built")
	cmd.PersistentFlags().String("version", defaultVersion, "project version being built")

	return nil
}
