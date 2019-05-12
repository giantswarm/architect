package triggerbuild

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/jszwedko/go-circleci"
	"github.com/spf13/cobra"
)

func runTriggerBuildError(cmd *cobra.Command, args []string) error {
	var (
		branch  = cmd.Flag("branch").Value.String()
		job     = cmd.Flag("job").Value.String()
		org     = cmd.Flag("org").Value.String()
		project = cmd.Flag("project").Value.String()
		token   = cmd.Flag("token").Value.String()
	)

	if org == "" {
		return microerror.Mask(missingOrganisationError)
	}

	if project == "" {
		return microerror.Mask(missingProjectError)
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
	build, err := client.ParameterizedBuild(org, project, branch, params)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Printf("build triggered: %s\n", build.BuildURL)

	return nil
}
