package events

import (
	"context"
	"log"

	"github.com/google/go-github/github"

	microerror "github.com/giantswarm/microkit/error"
)

// Environment is a name of an installation.
type Environment string

// environmentProjects is a mapping between Environments,
// and the projects that should be deployed there.
// We use codenames for installations to not expose customer information.
var environmentProjects = map[Environment][]string{
	"centaur": []string{
		"api",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"flannel-operator",
		"g8s-etcd-operator",
		"g8s-grafana",
		"g8s-prometheus",
		"happa",
		"ingress-operator",
		"kubernetesd",
		"kvm-operator",
		"passage",
		"testbot",
		"tokend",
		"userd",
	},
	"viking": []string{
		"api",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"g8s-etcd-operator",
		"g8s-grafana",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"aws-operator",
		"passage",
		"testbot",
		"tokend",
		"userd",
	},
	"heisenberg": []string{
		"api",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"g8s-etcd-operator",
		"g8s-grafana",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"aws-operator",
		"passage",
		"testbot",
		"tokend",
		"userd",
	},
	"asgard": []string{
		"api",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"g8s-etcd-operator",
		"g8s-grafana",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"aws-operator",
		"passage",
		"testbot",
		"tokend",
		"userd",
	},
	"iris": []string{
		"api",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"g8s-etcd-operator",
		"g8s-grafana",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"aws-operator",
		"passage",
		"testbot",
		"tokend",
		"userd",
	},
}

// GetEnvironments takes a project name, and returns a list of environments
// where this project should be deployed to.
func GetEnvironments(project string) []Environment {
	environments := []Environment{}

	for environment, projects := range environmentProjects {
		for _, p := range projects {
			if project == p {
				environments = append(environments, environment)
			}
		}
	}

	return environments
}

func CreateDeploymentEvent(client *github.Client, environment Environment, organisation, project, sha string) error {
	falseBool := false
	environmentString := string(environment)
	requiredContexts := []string{}

	deploymentRequest := github.DeploymentRequest{
		Ref:              &sha,
		AutoMerge:        &falseBool,
		Environment:      &environmentString,
		RequiredContexts: &requiredContexts,
	}

	_, _, err := client.Repositories.CreateDeployment(
		context.TODO(),
		organisation,
		project,
		&deploymentRequest,
	)
	if err != nil {
		return microerror.MaskAnyf(err, "could not create deployment event for %v", environment)
	}

	log.Printf("created deployment event for %v for %v for %v", project, environment, sha)

	return nil
}
