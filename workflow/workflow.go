package workflow

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"

	"github.com/spf13/afero"
)

var (
	GoTestCommandName  = "go-test"
	GoBuildCommandName = "go-build"

	DockerBuildCommandName      = "docker-build"
	DockerRunVersionCommandName = "docker-run-version"
	DockerRunHelpCommandName    = "docker-run-help"

	DockerLoginCommandName = "docker-login"
	DockerPushCommandName  = "docker-push"

	KubectlClusterInfoCommandName = "kubectl-cluster-info"
	KubectlApplyCommandName       = "kubectl-apply"
)

type Workflow []commands.Command

func (w Workflow) String() string {
	if len(w) == 0 {
		return "{}"
	}

	cmdStrings := []string{}
	for _, cmd := range w {
		cmdStrings = append(cmdStrings, "\t"+cmd.String()+"\n")
	}

	return fmt.Sprintf("{\n%v}", strings.Join(cmdStrings, ""))
}

type KubernetesCluster struct {
	ApiServer      string
	IngressTag     string
	CaPath         string
	CrtPath        string
	KeyPath        string
	KubectlVersion string
}

type ProjectInfo struct {
	WorkingDirectory string
	Organisation     string
	Project          string
	Sha              string

	Registry       string
	DockerEmail    string
	DockerUsername string
	DockerPassword string

	KubernetesResourcesDirectoryPath string
	KubernetesClusters               []KubernetesCluster

	Goos          string
	Goarch        string
	GolangImage   string
	GolangVersion string
}

func NewBuild(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	if projectInfo.WorkingDirectory == "" {
		return nil, fmt.Errorf("working directory cannot be empty")
	}
	if projectInfo.Organisation == "" {
		return nil, fmt.Errorf("organisation cannot be empty")
	}
	if projectInfo.Project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	mainGoExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "main.go"))
	if err != nil {
		return nil, err
	}
	if mainGoExists {
		if projectInfo.Goos == "" {
			return nil, fmt.Errorf("goos cannot be empty")
		}
		if projectInfo.Goarch == "" {
			return nil, fmt.Errorf("goarch cannot be empty")
		}
		if projectInfo.GolangImage == "" {
			return nil, fmt.Errorf("golang image cannot be empty")
		}
		if projectInfo.GolangVersion == "" {
			return nil, fmt.Errorf("golang version cannot be empty")
		}

		testPackageArguments, err := utils.NoVendor(fs, projectInfo.WorkingDirectory)
		if err != nil {
			return nil, err
		}

		goTest := commands.NewDockerCommand(
			GoTestCommandName,
			commands.DockerCommandConfig{
				Volumes: []string{
					fmt.Sprintf(
						"%v:/go/src/github.com/%v/%v",
						projectInfo.WorkingDirectory,
						projectInfo.Organisation,
						projectInfo.Project,
					),
				},
				Env: []string{
					fmt.Sprintf("GOOS=%v", projectInfo.Goos),
					fmt.Sprintf("GOARCH=%v", projectInfo.Goarch),
					"GOPATH=/go",
					"CGOENABLED=0",
				},
				WorkingDirectory: fmt.Sprintf(
					"/go/src/github.com/%v/%v",
					projectInfo.Organisation,
					projectInfo.Project,
				),
				Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
				Args:  []string{"go", "test", "-v"},
			},
		)
		goTest.Args = append(goTest.Args, testPackageArguments...)
		w = append(w, goTest)

		goBuild := commands.NewDockerCommand(
			GoBuildCommandName,
			commands.DockerCommandConfig{
				Volumes: []string{
					fmt.Sprintf(
						"%v:/go/src/github.com/%v/%v",
						projectInfo.WorkingDirectory,
						projectInfo.Organisation,
						projectInfo.Project,
					),
				},
				Env: []string{
					fmt.Sprintf("GOOS=%v", projectInfo.Goos),
					fmt.Sprintf("GOARCH=%v", projectInfo.Goarch),
					"GOPATH=/go",
					"CGOENABLED=0",
				},
				WorkingDirectory: fmt.Sprintf(
					"/go/src/github.com/%v/%v",
					projectInfo.Organisation,
					projectInfo.Project,
				),
				Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
				Args:  []string{"go", "build", "-v", "-a", "-tags", "netgo", "-ldflags", "-linkmode 'external' -extldflags '-static'"},
			},
		)
		w = append(w, goBuild)
	}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	if dockerFileExists {
		if projectInfo.Registry == "" {
			return nil, fmt.Errorf("registry cannot be empty")
		}
		if projectInfo.Sha == "" {
			return nil, fmt.Errorf("sha cannot be empty")
		}

		dockerBuild := commands.Command{
			Name: DockerBuildCommandName,
			Args: []string{
				"docker",
				"build",
				"-t",
				fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
				projectInfo.WorkingDirectory,
			},
		}
		w = append(w, dockerBuild)

		dockerRunVersion := commands.NewDockerCommand(
			DockerRunVersionCommandName,
			commands.DockerCommandConfig{
				Image: fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
				Args:  []string{"version"},
			},
		)
		w = append(w, dockerRunVersion)

		dockerRunHelp := commands.NewDockerCommand(
			DockerRunHelpCommandName,
			commands.DockerCommandConfig{
				Image: fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
				Args:  []string{"--help"},
			},
		)
		w = append(w, dockerRunHelp)
	}

	return w, nil
}

func NewDeploy(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	if projectInfo.WorkingDirectory == "" {
		return nil, fmt.Errorf("working directory cannot be empty")
	}
	if projectInfo.Organisation == "" {
		return nil, fmt.Errorf("organisation cannot be empty")
	}
	if projectInfo.Project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	if dockerFileExists {
		if projectInfo.Sha == "" {
			return nil, fmt.Errorf("sha cannot be empty")
		}
		if projectInfo.Registry == "" {
			return nil, fmt.Errorf("registry cannot be empty")
		}
		if projectInfo.DockerEmail == "" {
			return nil, fmt.Errorf("docker email cannot be empty")
		}
		if projectInfo.DockerUsername == "" {
			return nil, fmt.Errorf("docker username cannot be empty")
		}
		if projectInfo.DockerPassword == "" {
			return nil, fmt.Errorf("docker password cannot be empty")
		}

		dockerLogin := commands.Command{
			Name: DockerLoginCommandName,
			Args: []string{
				"docker",
				"login",
				fmt.Sprintf("--email=%v", projectInfo.DockerEmail),
				fmt.Sprintf("--username=%v", projectInfo.DockerUsername),
				fmt.Sprintf("--password=%v", projectInfo.DockerPassword),
				projectInfo.Registry,
			},
		}
		w = append(w, dockerLogin)

		dockerPush := commands.Command{
			Name: DockerPushCommandName,
			Args: []string{
				"docker",
				"push",
				fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
			},
		}
		w = append(w, dockerPush)
	}

	kubernetesDirectoryExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "kubernetes/"))
	if err != nil {
		return nil, err
	}
	if kubernetesDirectoryExists {
		if projectInfo.KubernetesResourcesDirectoryPath == "" {
			return nil, fmt.Errorf("kubernetes templated resources directory path cannot be empty")
		}

		for _, cluster := range projectInfo.KubernetesClusters {
			if cluster.IngressTag == "" {
				return nil, fmt.Errorf("ingress tag cannot be empty")
			}

			// Copy /kubernetes to a per-cluster directory, and template it
			dir, subdir := filepath.Split(projectInfo.KubernetesResourcesDirectoryPath)
			templatedResourcesDirectory := filepath.Join(dir, subdir+"-"+cluster.IngressTag)

			if err := utils.CopyDir(
				fs,
				projectInfo.KubernetesResourcesDirectoryPath,
				templatedResourcesDirectory,
			); err != nil {
				return nil, err
			}

			if err := utils.TemplateKubernetesResources(fs, templatedResourcesDirectory, projectInfo.Sha, cluster.IngressTag); err != nil {
				return nil, err
			}

			if cluster.ApiServer == "" {
				return nil, fmt.Errorf("kubernetes api server cannot be empty")
			}
			if cluster.CaPath == "" {
				return nil, fmt.Errorf("kubernetes ca path cannot be empty")
			}
			if cluster.CrtPath == "" {
				return nil, fmt.Errorf("kubernetes crt path cannot be empty")
			}
			if cluster.KeyPath == "" {
				return nil, fmt.Errorf("kubernetes key path cannot be empty")
			}
			if cluster.KubectlVersion == "" {
				return nil, fmt.Errorf("kubectl version cannot be empty")
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
			w = append(w, kubectlClusterInfo)

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
			w = append(w, kubectlApply)
		}
	}

	return w, nil
}

func getCertsFromEnv(fs afero.Fs, workingDirectory, envVarPrefix string) (string, string, string, error) {
	certDetails := []struct {
		envVarSuffix   string
		fileNameSuffix string
	}{
		{envVarSuffix: "_CA", fileNameSuffix: "-ca.pem"},
		{envVarSuffix: "_CRT", fileNameSuffix: "-crt.pem"},
		{envVarSuffix: "_KEY", fileNameSuffix: "-key.pem"},
	}

	filePaths := []string{}

	for _, certDetail := range certDetails {
		envVarName := envVarPrefix + certDetail.envVarSuffix

		certData := os.Getenv(envVarName)
		if certData == "" {
			return "", "", "", fmt.Errorf("could not find cert var: %v", envVarName)
		}

		certFileData, err := base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return "", "", "", fmt.Errorf("could not decode cert: %v", err)
		}

		filePath := filepath.Join(
			workingDirectory,
			strings.ToLower(envVarPrefix)+certDetail.fileNameSuffix,
		)
		if err := afero.WriteFile(fs, filePath, certFileData, 0644); err != nil {
			return "", "", "", fmt.Errorf("could not write cert: %v", err)
		}

		filePaths = append(filePaths, filePath)
	}

	if len(filePaths) != 3 {
		return "", "", "", fmt.Errorf("incorrect number of certs found")
	}

	return filePaths[0], filePaths[1], filePaths[2], nil
}

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
		caPath, crtPath, keyPath, err := getCertsFromEnv(fs, workingDirectory, configuredCluster.EnvVarPrefix)
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
