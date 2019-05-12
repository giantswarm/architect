package triggerjob

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/jszwedko/go-circleci"
	"github.com/spf13/cobra"
)

func runTriggerJobError(cmd *cobra.Command, args []string) error {
	var (
		branch = cmd.Flag("branch").Value.String()
		job    = cmd.Flag("job").Value.String()
		org    = cmd.Flag("org").Value.String()
		repo   = cmd.Flag("repo").Value.String()
		token  = cmd.Flag("token").Value.String()
	)

	if org == "" {
		return microerror.Mask(missingOrganisationError)
	}

	if repo == "" {
		return microerror.Mask(missingRepositoryError)
	}

	if branch == "" {
		return microerror.Mask(missingBranchError)
	}

	var params map[string]string = nil
	if job != "" {
		params = map[string]string{
			"CIRCLE_JOB": job,
		}
	}

	client := &circleci.Client{Token: token}
	build, err := client.ParameterizedBuild(org, repo, branch, params)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Printf("job triggered: %s\n", build.BuildURL)

	return nil
}
