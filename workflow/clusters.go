package workflow

import "github.com/spf13/afero"

func ClustersFromEnv(fs afero.Fs, workingDirectory string) ([]KubernetesCluster, error) {
	type Cluster struct {
		ApiServer      string
		IngressTag     string
		EnvVarPrefix   string
		KubectlVersion string
	}
	configuredClusters := []Cluster{
		Cluster{
			ApiServer:      "https://api.g8s.fra-1.giantswarm.io",
			IngressTag:     "g8s",
			EnvVarPrefix:   "G8S",
			KubectlVersion: "1.4.7",
		},
		Cluster{
			ApiServer:      "https://api.g8s.eu-west-1.aws.adidas.private.giantswarm.io:6443",
			IngressTag:     "aws",
			EnvVarPrefix:   "AWS",
			KubectlVersion: "1.4.7",
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
			IngressTag:     configuredCluster.IngressTag,
			CaPath:         caPath,
			CrtPath:        crtPath,
			KeyPath:        keyPath,
			KubectlVersion: configuredCluster.KubectlVersion,
		}
		clusters = append(clusters, newCluster)
	}

	return clusters, nil
}
