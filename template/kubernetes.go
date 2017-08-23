package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/utils"
)

const (
	TemplateKubernetesResourcesTaskName = "template-kubernetes-resources"

	// TemplateKubernetesResourcesTaskString is the format for printing the
	// kubernetes resources templating task.
	// Name of the task, templated resources directory path, sha.
	TemplateKubernetesResourcesTaskString = "%s:\t%s %s"
)

type TemplateKubernetesResourcesTask struct {
	fs afero.Fs

	kubernetesResourcesDirectoryPath string
	templatedResourcesDirectoryPath  string

	sha          string
	installation configuration.Installation
}

func (t TemplateKubernetesResourcesTask) Run() error {
	if t.kubernetesResourcesDirectoryPath == "" {
		return microerror.Mask(emptyKubernetesResourcesDirectoryPath)
	}

	if err := utils.CopyDir(
		t.fs,
		t.kubernetesResourcesDirectoryPath,
		t.templatedResourcesDirectoryPath,
	); err != nil {
		return microerror.Mask(err)
	}

	config := TemplateConfiguration{
		BuildInfo: BuildInfo{
			SHA: t.sha,
		},
		Installation: t.installation,
	}

	if err := templateKubernetesResources(t.fs, t.templatedResourcesDirectoryPath, config); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (t TemplateKubernetesResourcesTask) Name() string {
	return TemplateKubernetesResourcesTaskName
}

func (t TemplateKubernetesResourcesTask) String() string {
	return fmt.Sprintf(
		TemplateKubernetesResourcesTaskString,
		t.Name(),
		t.templatedResourcesDirectoryPath,
		t.sha,
	)
}

func NewTemplateKubernetesResourcesTask(fs afero.Fs, kubernetesResourcesDirectoryPath, templatedResourcesDirectoryPath, sha string, installation configuration.Installation) TemplateKubernetesResourcesTask {
	return TemplateKubernetesResourcesTask{
		fs: fs,

		kubernetesResourcesDirectoryPath: kubernetesResourcesDirectoryPath,
		templatedResourcesDirectoryPath:  templatedResourcesDirectoryPath,

		sha:          sha,
		installation: installation,
	}
}

// templateKubernetesResources takes a filesystem,
// a path to a directory holding kubernetes resources,
// and an installation configuration.
// It templates the given resources, with data from the configuration,
// writing changes to the files.
func templateKubernetesResources(fs afero.Fs, resourcesPath string, config TemplateConfiguration) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			microerror.Mask(err)
		}

		if info.IsDir() {
			return nil
		}

		contents, err := afero.ReadFile(fs, path)
		if err != nil {
			microerror.Mask(err)
		}

		t := template.Must(template.New(path).Funcs(filters).Parse(string(contents)))
		if err != nil {
			microerror.Mask(err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, config); err != nil {
			microerror.Mask(err)
		}

		templatedContents := buf.String()

		// This adds backwards compatability for `%%DOCKER_TAG%%`. Deprecated.
		templatedContents = strings.Replace(templatedContents, "%%DOCKER_TAG%%", config.BuildInfo.SHA, -1)

		if err := afero.WriteFile(fs, path, []byte(templatedContents), permission); err != nil {
			microerror.Mask(err)
		}

		return nil
	}

	if err := afero.Walk(fs, resourcesPath, walkFunc); err != nil {
		microerror.Mask(err)
	}

	return nil
}
