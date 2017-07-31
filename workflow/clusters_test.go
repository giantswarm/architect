package workflow

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestClustersFromEnv(t *testing.T) {
	type envVar struct {
		key   string
		value string
	}

	certTestData := "test"
	certEncodedTestData := "dGVzdA=="

	tests := []struct {
		workingDirectory string
		envVars          []envVar
		clusters         []KubernetesCluster
		expectedFiles    []string
	}{
		// Test that an empty env var set returns no clusters.
		{
			workingDirectory: "/test/",
			envVars:          []envVar{},
			clusters:         []KubernetesCluster{},
			expectedFiles:    []string{},
		},

		// Test that a few G8S env vars returns no clusters.
		{
			workingDirectory: "/test/",
			envVars: []envVar{
				envVar{key: "G8S_CA", value: "test"},
			},
			clusters:      []KubernetesCluster{},
			expectedFiles: []string{},
		},

		// Test that G8S env vars returns the LW cluster
		{
			workingDirectory: "/test/",
			envVars: []envVar{
				envVar{key: "G8S_CA", value: certEncodedTestData},
				envVar{key: "G8S_CRT", value: certEncodedTestData},
				envVar{key: "G8S_KEY", value: certEncodedTestData},
			},
			clusters: []KubernetesCluster{
				KubernetesCluster{
					ApiServer:      "https://api.g8s.fra-1.giantswarm.io",
					Prefix:         "G8S",
					CaPath:         "/test/g8s-ca.pem",
					CrtPath:        "/test/g8s-crt.pem",
					KeyPath:        "/test/g8s-key.pem",
					KubectlVersion: "f51f93c30d27927d2b33122994c0929b3e6f2432",
				},
			},
			expectedFiles: []string{
				"/test/g8s-ca.pem",
				"/test/g8s-crt.pem",
				"/test/g8s-key.pem",
			},
		},

		// Test that AWS env vars returns the AWS cluster
		{
			workingDirectory: "/test",
			envVars: []envVar{
				envVar{key: "AWS_CA", value: certEncodedTestData},
				envVar{key: "AWS_CRT", value: certEncodedTestData},
				envVar{key: "AWS_KEY", value: certEncodedTestData},
			},
			clusters: []KubernetesCluster{
				KubernetesCluster{
					ApiServer:      "https://api.g8s.eu-west-1.aws.adidas.private.giantswarm.io:6443",
					Prefix:         "AWS",
					CaPath:         "/test/aws-ca.pem",
					CrtPath:        "/test/aws-crt.pem",
					KeyPath:        "/test/aws-key.pem",
					KubectlVersion: "a121f8d14cd14567abc2ec20a7258be9d70ecb45",
				},
			},
			expectedFiles: []string{
				"/test/aws-ca.pem",
				"/test/aws-crt.pem",
				"/test/aws-key.pem",
			},
		},

		// Test that G8S and AWS env vars return both clusters
		{
			workingDirectory: "/test",
			envVars: []envVar{
				envVar{key: "G8S_CA", value: certEncodedTestData},
				envVar{key: "G8S_CRT", value: certEncodedTestData},
				envVar{key: "G8S_KEY", value: certEncodedTestData},

				envVar{key: "AWS_CA", value: certEncodedTestData},
				envVar{key: "AWS_CRT", value: certEncodedTestData},
				envVar{key: "AWS_KEY", value: certEncodedTestData},
			},
			clusters: []KubernetesCluster{
				KubernetesCluster{
					ApiServer:      "https://api.g8s.fra-1.giantswarm.io",
					Prefix:         "G8S",
					CaPath:         "/test/g8s-ca.pem",
					CrtPath:        "/test/g8s-crt.pem",
					KeyPath:        "/test/g8s-key.pem",
					KubectlVersion: "f51f93c30d27927d2b33122994c0929b3e6f2432",
				},
				KubernetesCluster{
					ApiServer:      "https://api.g8s.eu-west-1.aws.adidas.private.giantswarm.io:6443",
					Prefix:         "AWS",
					CaPath:         "/test/aws-ca.pem",
					CrtPath:        "/test/aws-crt.pem",
					KeyPath:        "/test/aws-key.pem",
					KubectlVersion: "a121f8d14cd14567abc2ec20a7258be9d70ecb45",
				},
			},
			expectedFiles: []string{
				"/test/g8s-ca.pem",
				"/test/g8s-crt.pem",
				"/test/g8s-key.pem",

				"/test/aws-ca.pem",
				"/test/aws-crt.pem",
				"/test/aws-key.pem",
			},
		},
	}

	for index, test := range tests {
		for _, envVar := range test.envVars {
			if err := os.Setenv(envVar.key, envVar.value); err != nil {
				t.Fatalf("%v: unexpected error setting env var: %v", index, err)
			}
		}

		fs := afero.NewMemMapFs()

		clusters, err := ClustersFromEnv(fs, test.workingDirectory)

		for _, envVar := range test.envVars {
			if err := os.Setenv(envVar.key, ""); err != nil {
				t.Fatalf("%v: unexpected error unsetting env var: %v", index, err)
			}
		}

		if err != nil {
			t.Fatalf("%v: unexpected error getting clusters from env: %v", index, err)
		}

		for index := range clusters {
			if clusters[index].ApiServer != test.clusters[index].ApiServer {
				t.Fatalf(
					"%v: expected api server did not match returned api server.\nexpected: %#v\nreturned: %#v\n",
					index,
					test.clusters[index].ApiServer,
					clusters[index].ApiServer,
				)
			}

			if clusters[index].Prefix != test.clusters[index].Prefix {
				t.Fatalf(
					"%v: expected prefix did not match returned prefix.\nexpected: %#v\nreturned: %#v\n",
					index,
					test.clusters[index].Prefix,
					clusters[index].Prefix,
				)
			}

			if clusters[index].KubectlVersion != test.clusters[index].KubectlVersion {
				t.Fatalf(
					"%v: expected kubectl version did not match returned kubectl version.\nexpected: %#v\nreturned: %#v\n",
					index,
					test.clusters[index].KubectlVersion,
					clusters[index].KubectlVersion,
				)
			}
		}

		for _, expectedFile := range test.expectedFiles {
			_, err := fs.Stat(expectedFile)
			if err != nil {
				t.Fatalf("%v: unexpected error checking certificate: %v", index, err)
			}

			contents, err := afero.ReadFile(fs, expectedFile)
			if err != nil {
				t.Fatalf("%v: unexpected error checking certificate contents: %v", index, err)
			}
			if string(contents) != certTestData {
				t.Fatalf("%v: certificate did not match expected contents: %v", index, string(contents))
			}
		}
	}
}
