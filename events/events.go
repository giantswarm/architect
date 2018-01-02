package events

import (
	"context"
	"log"

	"github.com/google/go-github/github"

	"github.com/giantswarm/microerror"
)

// Environment is a name of an installation.
type Environment string

// environmentProjects is a mapping between Environments,
// and the projects that should be deployed there.
// We use codenames for installations to not expose customer information.
var environmentProjects = map[Environment][]string{
	"centaur": {
		"api",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"endpoint-operator",
		"flannel-operator",
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
	"viking": {
		"api",
		"calico-node-controller",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"asgard": {
		"api",
		"calico-node-controller",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"iris": {
		"api",
		"calico-node-controller",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"anubis": {
		"api",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"endpoint-operator",
		"flannel-operator",
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
	"ginger": {
		"api",
		"cert-exporter",
		"calico-node-controller",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"ginger_vxvc7": {
		"api-spec",
		"docs",
		"docs-indexer",
		"docs-proxy",
		"giantswarmio-nginx",
		"giantswarmio-webapp",
		"web-assets",
		"sitesearch",
		"classify",
		"mole-rat",
	},
	"gauss": {
		"api",
		"cert-exporter",
		"calico-node-controller",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"archon": {
		"api",
		"calico-node-controller",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"jackal": {
		"api",
		"calico-node-controller",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"avatar": {
		"api",
		"calico-node-controller",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-cloudwatch-exporter",
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
	"amagon": {
		"api",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"endpoint-operator",
		"flannel-operator",
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
	"geckon": {
		"api",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"endpoint-operator",
		"flannel-operator",
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
	"gollum": {
		"api",
		"azure-operator",
		"cert-exporter",
		"cert-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"draughtsman-eventer",
		"draughtsman-operator",
		"etcd-backup",
		"g8s-grafana",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
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
		return microerror.Maskf(err, "could not create deployment event for %v", environment)
	}

	log.Printf("created deployment event for %v for %v for %v", project, environment, sha)

	return nil
}
