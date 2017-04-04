package workflow

import (
	"fmt"

	"github.com/giantswarm/architect/commands"
	"github.com/spf13/afero"
)

var (
	KubectlClusterInfoCommandName = "kubectl-cluster-info"
	KubectlApplyCommandName       = "kubectl-apply"
)

func checkKubectlRequirements(cluster KubernetesCluster) error {
	if cluster.ApiServer == "" {
		return fmt.Errorf("kubernetes api server cannot be empty")
	}
	if cluster.CaPath == "" {
		return fmt.Errorf("kubernetes ca path cannot be empty")
	}
	if cluster.CrtPath == "" {
		return fmt.Errorf("kubernetes crt path cannot be empty")
	}
	if cluster.KeyPath == "" {
		return fmt.Errorf("kubernetes key path cannot be empty")
	}
	if cluster.KubectlVersion == "" {
		return fmt.Errorf("kubectl version cannot be empty")
	}

	return nil
}

func NewKubectlClusterInfoCommand(fs afero.Fs, cluster KubernetesCluster) (commands.Command, error) {
	if err := checkKubectlRequirements(cluster); err != nil {
		return commands.Command{}, err
	}

	kubectlClusterInfo := commands.NewDockerCommand(
		KubectlClusterInfoCommandName,
		commands.DockerCommandConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
				fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
				fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
			},
			Image: fmt.Sprintf("giantswarm/kubectl:%v", cluster.KubectlVersion),
			Args: []string{
				fmt.Sprintf("--server=%v", cluster.ApiServer),
				"--certificate-authority=/ca.pem",
				"--client-certificate=/crt.pem",
				"--client-key=/key.pem",
				"cluster-info",
			},
		},
	)

	return kubectlClusterInfo, nil
}

func NewKubectlApplyCommand(fs afero.Fs, cluster KubernetesCluster, templatedResourcesDirectory string) (commands.Command, error) {
	if err := checkKubectlRequirements(cluster); err != nil {
		return commands.Command{}, err
	}

	kubectlApply := commands.NewDockerCommand(
		KubectlApplyCommandName,
		commands.DockerCommandConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
				fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
				fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
				fmt.Sprintf("%v:/kubernetes", templatedResourcesDirectory),
			},
			Image: fmt.Sprintf("giantswarm/kubectl:%v", cluster.KubectlVersion),
			Args: []string{
				fmt.Sprintf("--server=%v", cluster.ApiServer),
				"--certificate-authority=/ca.pem",
				"--client-certificate=/crt.pem",
				"--client-key=/key.pem",
				"apply", "-R", "-f", "/kubernetes",
			},
		},
	)

	return kubectlApply, nil
}
