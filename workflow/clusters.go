package workflow

import (
	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/installation"
	"github.com/spf13/afero"
)

type KubernetesCluster struct {
	ApiServer      string
	Prefix         string
	CaPath         string
	CrtPath        string
	KeyPath        string
	KubectlVersion string

	Installation configuration.Installation
}

func ClustersFromEnv(fs afero.Fs, workingDirectory string) ([]KubernetesCluster, error) {
	type Cluster struct {
		ApiServer      string
		EnvVarPrefix   string
		KubectlVersion string

		Installation configuration.Installation
	}
	configuredClusters := []Cluster{
		Cluster{
			ApiServer:      "https://api.g8s.fra-1.giantswarm.io",
			EnvVarPrefix:   "G8S",
			KubectlVersion: "f9a5a46a94cd951d631a9d8fc0e95c6c73a46fb2",

			Installation: installation.Leaseweb,
		},
		Cluster{
			ApiServer:      "https://api.g8s.eu-west-1.aws.adidas.private.giantswarm.io:6443",
			EnvVarPrefix:   "AWS",
			KubectlVersion: "f9a5a46a94cd951d631a9d8fc0e95c6c73a46fb2",

			Installation: installation.AWS,
		},
	}

	clusters := []KubernetesCluster{}

	for _, configuredCluster := range configuredClusters {
		caPath, crtPath, keyPath, err := CertsFromEnv(fs, workingDirectory, configuredCluster.EnvVarPrefix)
		if err != nil {
			continue
		}

		newCluster := KubernetesCluster{
			ApiServer:      configuredCluster.ApiServer,
			Prefix:         configuredCluster.EnvVarPrefix,
			CaPath:         caPath,
			CrtPath:        crtPath,
			KeyPath:        keyPath,
			KubectlVersion: configuredCluster.KubectlVersion,

			Installation: configuredCluster.Installation,
		}
		clusters = append(clusters, newCluster)
	}

	return clusters, nil
}
