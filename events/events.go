package events

import (
	"context"
	"log"

	"github.com/google/go-github/github"

	"github.com/giantswarm/microerror"
)

var (
	awsProjectList = []string{
		"api",
		"cert-exporter",
		"cert-operator",
		"cluster-operator",
		"cluster-service",
		"companyd",
		"credentiald",
		"desmotes",
		"draughtsman",
		"etcd-backup",
		"g8s-cert-manager",
		"g8s-cloudwatch-exporter",
		"g8s-efk",
		"g8s-grafana",
		"g8s-oauth2-proxy",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"aws-operator",
		"node-operator",
		"passage",
		"tokend",
		"userd",
		"vault-exporter",
	}
	azureProjectList = []string{
		"api",
		"azure-operator",
		"cert-exporter",
		"cert-operator",
		"cluster-operator",
		"cluster-service",
		"companyd",
		"credentiald",
		"desmotes",
		"draughtsman",
		"etcd-backup",
		"g8s-cert-manager",
		"g8s-efk",
		"g8s-grafana",
		"g8s-oauth2-proxy",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"node-operator",
		"passage",
		"tokend",
		"userd",
		"vault-exporter",
	}
	kvmProjectList = []string{
		"api",
		"cert-exporter",
		"cert-operator",
		"cluster-operator",
		"cluster-service",
		"companyd",
		"desmotes",
		"draughtsman",
		"etcd-backup",
		"endpoint-operator",
		"flannel-operator",
		"g8s-cert-manager",
		"g8s-efk",
		"g8s-grafana",
		"g8s-oauth2-proxy",
		"g8s-prometheus",
		"happa",
		"ingress-operator",
		"kubernetesd",
		"kvm-operator",
		"node-operator",
		"passage",
		"tokend",
		"userd",
		"vault-exporter",
	}
)

// Environment is a name of an installation.
type Environment string

// environmentProjects is a mapping between Environments,
// and the projects that should be deployed there.
// We use codenames for installations to not expose customer information.
var environmentProjects = map[Environment][]string{
	"centaur": append(kvmProjectList, "pv-cleaner-operator"),
	"viking":  awsProjectList,
	"asgard":  awsProjectList,
	"iris":    awsProjectList,
	"anubis":  kvmProjectList,
	"ginger":  awsProjectList,
	"gauss":   awsProjectList,
	"archon":  awsProjectList,
	"jackal":  awsProjectList,
	"avatar":  awsProjectList,
	"amagon":  kvmProjectList,
	"geckon":  kvmProjectList,
	"gollum":  azureProjectList,
	"gollum_6f342": {
		"api-spec",
		"docs",
		"docs-indexer",
		"docs-proxy",
		"draughtsman",
		"giantswarmio-nginx",
		"giantswarmio-webapp",
		"web-assets",
		"sitesearch",
	},
	"gorgoth":   kvmProjectList,
	"atlantis":  awsProjectList,
	"godsmack":  azureProjectList,
	"victory":   azureProjectList,
	"tarantula": azureProjectList,
	"axolotl":   awsProjectList,
	"ghost":     azureProjectList,
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
