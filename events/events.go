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
		"app-operator",
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
		"operator-week-app-collection",
		"passage",
		"release-operator",
		"tokend",
		"userd",
		"vault-exporter",
	}
	awsProjectList = append(baseProjectList,
		"aws-operator",
		"aws-app-collection",
	)
	// azure project list + azure specific services
	azureProjectList = append(baseProjectList,
		"azure-operator",
		"azure-app-collection",
	)
	// kvm project list + kvm specific services
	kvmProjectList = append(baseProjectList,
		"bridge-operator",
		"flannel-operator",
		"ingress-exporter",
		"kvm-operator",
		"kvm-app-collection",
	)
	// awsChinaProjectList adds route53-manager to awsProjectList. This is
	// required until route53 is available on AWS China.
	awsChinaProjectList = append(awsProjectList,
		"route53-manager",
	)
)

// Environment is a name of an installation.
type Environment string

// environmentProjects is a mapping between Environments,
// and the projects that should be deployed there.
// We use codenames for installations to not expose customer information.
var environmentProjects = map[Environment][]string{
	"amagon":   kvmProjectList,
	"anubis":   kvmProjectList,
	"archon":   awsProjectList,
	"asgard":   awsProjectList,
	"atlantis": awsProjectList,
	"avatar":   awsProjectList,
	"axolotl":  awsChinaProjectList,
	"buffalo":  kvmProjectList,
	"centaur":  append(kvmProjectList, "pv-cleaner-operator"),
	"davis":    azureProjectList,
	"dinosaur": kvmProjectList,
	"dragon":   kvmProjectList,
	"gauss":    awsProjectList,
	"geckon":   kvmProjectList,
	"ghost":    azureProjectList,
	"ginger":   awsProjectList,
	"giraffe":  awsChinaProjectList,
	"godsmack": azureProjectList,
	"gollum":   azureProjectList,
	"goku":     azureProjectList,
	"gorgoth":  kvmProjectList,
	"gorilla":  awsProjectList,
	"iris":     awsProjectList,
	"puma":     kvmProjectList,
	"panther":  kvmProjectList,
	"seal":     awsProjectList,
	"talos":    azureProjectList,
	"victory":  azureProjectList,
	"viking":   awsProjectList,

	// non standard cluster just for our website
	"gollum_6iec4": {
		"auto-oncall",
		"circleci-exporter",
		"draughtsman",
		"github-exporter",
		"kpi-prometheus",
	},
}

// testingGroup is a specific grouping of environments that are considered
// testing environments.
var testingGroup = []Environment{
	"gauss",
	"geckon",
	"ghost",
	"ginger",
	"giraffe",
	"godsmack",
	"gorgoth",
}

// isTestingEnvironment checks if a given environment is in the testingGroup.
func isTestingEnvironment(environment Environment) bool {
	for _, testEnvironment := range testingGroup {
		if environment == testEnvironment {
			return true
		}
	}

	return false
}

// GetEnvironments takes a project name, and returns a list of environments
// where this project should be deployed to. If group is 'testing', then it
// only considers environments in the testingGroup.
func GetEnvironments(project string, group string) []Environment {
	environments := []Environment{}

	for environment, projects := range environmentProjects {
		// Skip environments that are not testing environments
		// if we are requesting deploys to the testing group
		if group == "testing" && !isTestingEnvironment(environment) {
			continue
		}

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
