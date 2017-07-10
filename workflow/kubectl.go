package workflow

import (
	"fmt"

	"github.com/giantswarm/architect/tasks"
	"github.com/spf13/afero"
)

var (
	KubectlClusterInfoTaskName = "kubectl-cluster-info"
	KubectlApplyTaskName       = "kubectl-apply"
)

func checkKubectlRequirements(cluster KubernetesCluster) error {
	if cluster.ApiServer == "" {
		return emptyKubernetesAPIServerError
	}
	if cluster.CaPath == "" {
		return emptyKubernetesCaPathError
	}
	if cluster.CrtPath == "" {
		return emptyKubernetesCrtPathError
	}
	if cluster.KeyPath == "" {
		return emptyKubernetesKeyPathError
	}
	if cluster.KubectlVersion == "" {
		return emptyKubectlVersionError
	}

	return nil
}

func NewKubectlClusterInfoTask(fs afero.Fs, cluster KubernetesCluster) (tasks.Task, error) {
	if err := checkKubectlRequirements(cluster); err != nil {
		return nil, err
	}

	kubectlClusterInfo := tasks.NewDockerTask(
		KubectlClusterInfoTaskName,
		tasks.DockerTaskConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
				fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
				fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
			},
			Image: fmt.Sprintf("quay.io/giantswarm/docker-kubectl:%v", cluster.KubectlVersion),
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

func NewKubectlApplyTask(fs afero.Fs, cluster KubernetesCluster, templatedResourcesDirectory string) (tasks.Task, error) {
	if err := checkKubectlRequirements(cluster); err != nil {
		return nil, err
	}

	kubectlApply := tasks.NewDockerTask(
		KubectlApplyTaskName,
		tasks.DockerTaskConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
				fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
				fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
				fmt.Sprintf("%v:/kubernetes", templatedResourcesDirectory),
			},
			Image: fmt.Sprintf("quay.io/giantswarm/docker-kubectl:%v", cluster.KubectlVersion),
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
