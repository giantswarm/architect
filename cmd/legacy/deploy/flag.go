package deploy

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagDeploymentEventsToken = "deployment-events-token"
	flagOrganisation          = "organisation"
	flagProject               = "project"
)

type flag struct {
	DeploymentEventsToken string
	Organisation          string
	Project               string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.DeploymentEventsToken, flagDeploymentEventsToken, "", `GitHub token used to create GitHub Deployment events. Defaults to the value of DEPLOYMENT_EVENTS_TOKEN environment variable.`)
}

func (f *flag) Validate() error {
	if f.DeploymentEventsToken == "" {
		f.DeploymentEventsToken = os.Getenv("DEPLOYMENT_EVENTS_TOKEN")
	}
	if f.DeploymentEventsToken == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagDeploymentEventsToken)
	}

	return nil
}
