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
		"cluster-operator",
		"draughtsman",
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
		"kvm-operator",
		"kvm-app-collection",
	)
)

// Environment is a name of an installation.
type Environment string

// environmentProjects is a mapping between Environments,
// and the projects that should be deployed there.
// We use codenames for installations to not expose customer information.
var environmentProjects = map[Environment][]string{
	"alpaca":   awsProjectList,
	"amagon":   kvmProjectList,
	"anteater": awsProjectList,
	"antelope": awsProjectList,
	"anubis":   kvmProjectList,
	"archon":   awsProjectList,
	"argali":   awsProjectList,
	"asgard":   awsProjectList,
	"atlantis": awsProjectList,
	"avatar":   awsProjectList,
	"axolotl":  awsProjectList,
	"buffalo":  kvmProjectList,
	"centaur":  append(kvmProjectList, "pv-cleaner-operator"),
	"davis":    azureProjectList,
	"dinosaur": kvmProjectList,
	"dragon":   kvmProjectList,
	"gaia":     awsProjectList,
	"gauss":    awsProjectList,
	"geckon":   kvmProjectList,
	"ghost":    azureProjectList,
	"ginger":   awsProjectList,
	"giraffe":  awsProjectList,
	"godsmack": azureProjectList,
	"goku":     azureProjectList,
	"gollum":   azureProjectList,
	"gorgoth":  kvmProjectList,
	"gorilla":  awsProjectList,
	"icarus":   awsProjectList,
	"impala":   azureProjectList,
	"iris":     awsProjectList,
	"panther":  kvmProjectList,
	"platypus": awsProjectList,
	"puma":     kvmProjectList,
	"seal":     awsProjectList,
	"talos":    azureProjectList,
	"victory":  azureProjectList,
	"viking":   awsProjectList,
}

// testingGroup is a specific grouping of environments that are considered
// testing environments.
var testingGroup = []Environment{
	"gaia",
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
		return microerror.Maskf(executionFailedError, "failed to create deployment event for %v with error %#q", environment, err)
	}

	log.Printf("created deployment event for %v for %v for %v", project, environment, sha)

	return nil
}
