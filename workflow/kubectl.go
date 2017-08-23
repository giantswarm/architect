package workflow

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/template"
)

var (
	KubectlClusterInfoTaskName = "kubectl-cluster-info"
	KubectlApplyTaskName       = "kubectl-apply"
)

func checkKubectlRequirements(cluster KubernetesCluster) error {
	if cluster.ApiServer == "" {
		return microerror.Mask(emptyKubernetesAPIServerError)
	}
	if cluster.CaPath == "" {
		return microerror.Mask(emptyKubernetesCaPathError)
	}
	if cluster.CrtPath == "" {
		return microerror.Mask(emptyKubernetesCrtPathError)
	}
	if cluster.KeyPath == "" {
		return microerror.Mask(emptyKubernetesKeyPathError)
	}
	if cluster.KubectlVersion == "" {
		return microerror.Mask(emptyKubectlVersionError)
	}

	return nil
}

func NewTemplateKubernetesResourcesTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkKubectlRequirements(projectInfo.CurrentCluster); err != nil {
		return nil, microerror.Mask(err)
	}

	templateKubernetesResources := template.NewTemplateKubernetesResourcesTask(
		fs,
		projectInfo.KubernetesResourcesDirectoryPath,
		projectInfo.TemplatedResourcesDirectory,
		projectInfo.Sha,
		projectInfo.CurrentCluster.Installation,
	)

	return templateKubernetesResources, nil
}

func NewKubectlClusterInfoTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	cluster := projectInfo.CurrentCluster

	if err := checkKubectlRequirements(cluster); err != nil {
		return nil, microerror.Mask(err)
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

func NewKubectlApplyTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	cluster := projectInfo.CurrentCluster

	if err := checkKubectlRequirements(cluster); err != nil {
		return nil, microerror.Mask(err)
	}

	kubectlApply := tasks.NewDockerTask(
		KubectlApplyTaskName,
		tasks.DockerTaskConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
				fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
				fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
				fmt.Sprintf("%v:/kubernetes", projectInfo.TemplatedResourcesDirectory),
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
