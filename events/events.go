package events

import (
	"context"
	"log"

	"github.com/google/go-github/github"

	"github.com/giantswarm/microerror"
)

var (
	// services used in all our installations
	baseProjectList = []string{
		"api",
		"cert-exporter",
		"cert-operator",
		"chart-operator",
		"cluster-operator",
		"cluster-service",
		"companyd",
		"credentiald",
		"draughtsman",
		"etcd-backup",
		"g8s-cert-manager",
		"g8s-efk",
		"g8s-grafana",
		"g8s-oauth2-proxy",
		"g8s-prometheus",
		"happa",
		"kubernetesd",
		"net-exporter",
		"node-operator",
		"passage",
		"release-operator",
		"tokend",
		"userd",
		"vault-exporter",
	}
	awsBaseProjectList = append(baseProjectList,
		"aws-operator",
	)
	// azure project list + azure specific services
	azureProjectList = append(baseProjectList,
		"azure-operator",
	)
	// kvm project list + kvm specific services
	kvmProjectList = append(baseProjectList,
		"endpoint-operator",
		"flannel-operator",
		"ingress-operator",
		"kvm-operator",
	)
	// awsChinaProjectList adds route53-manager to awsProjectList. This is
	// required until route53 is vailable on AWS China.
	awsChinaProjectList = append(awsBaseProjectList,
		"route53-manager",
	)
)

// Environment is a name of an installation.
type Environment string

// environmentProjects is a mapping between Environments,
// and the projects that should be deployed there.
// We use codenames for installations to not expose customer information.
var environmentProjects = map[Environment][]string{
	"amagon":    kvmProjectList,
	"anubis":    kvmProjectList,
	"archon":    awsBaseProjectList,
	"asgard":    awsBaseProjectList,
	"atlantis":  awsBaseProjectList,
	"avatar":    awsBaseProjectList,
	"axolotl":   awsChinaProjectList,
	"centaur":   append(kvmProjectList, "pv-cleaner-operator"),
	"dinosaur":  kvmProjectList,
	"dragon":    kvmProjectList,
	"gauss":     awsBaseProjectList,
	"geckon":    kvmProjectList,
	"ghost":     azureProjectList,
	"ginger":    awsBaseProjectList,
	"giraffe":   awsChinaProjectList,
	"godsmack":  azureProjectList,
	"gollum":    azureProjectList,
	"gorgoth":   kvmProjectList,
	"iris":      awsBaseProjectList,
	"seal":      awsBaseProjectList,
	"talos":     azureProjectList,
	"tarantula": azureProjectList,
	"victory":   azureProjectList,
	"viking":    awsBaseProjectList,

	// non standard cluster just for our website
	"gollum_6iec4": {
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
