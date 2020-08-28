package deploy

import (
	"context"
	"fmt"
	"io"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/giantswarm/architect/events"
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	// This should not be here but due to legacy structure those flags are
	// not initialised when calling Flag.Init.
	{
		r.flag.Organisation = cmd.Flag("organisation").Value.String()
		r.flag.Project = cmd.Flag("project").Value.String()

		if r.flag.Organisation == "" {
			return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagOrganisation)
		}
		if r.flag.Project == "" {
			return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagProject)
		}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.flag.DeploymentEventsToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	environments := events.GetEnvironments(r.flag.Project, events.GroupAll)

	fmt.Fprintf(r.stdout, "Creating deployment events for environments: %v\n", environments)
	for _, environment := range environments {
		err := events.CreateDeploymentEvent(githubClient, environment, r.flag.Organisation, r.flag.Project, cmd.Flag("sha").Value.String())
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
