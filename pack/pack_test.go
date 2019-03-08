package pack

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/microerror"
	"k8s.io/helm/pkg/chartutil"
	hapichart "k8s.io/helm/pkg/proto/hapi/chart"
)

var (
	chartYaml = `name: hello-test-chart
version: v1.0.0`
	valuesYaml     = `namespace: giantswarm`
	deploymentYaml = `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: hello-test
  namespace: {{ .Values.namespace }}
  labels:
    app: hello-test
spec:
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: hello-test
    spec:
      containers:
      - name: hello-test
        image: quay.io/giantswarm/hello-test:v1.0.0`
)

type File struct {
	path string
	data string
}

func TestPackageHelmChartTask(t *testing.T) {
	testCases := []struct {
		name string

		expectedChartName      string
		expectedChartVersion   string
		expectedDeploymentName string
		expectedDeploymentData string
		expectedValues         string

		errorMatcher func(err error) bool
		setup        func(chartDir string) error
	}{
		{
			name:                   "case 0: package chart",
			expectedChartName:      "hello-test-chart",
			expectedChartVersion:   "v1.0.0",
			expectedDeploymentName: "templates/deployment.yaml",
			expectedDeploymentData: deploymentYaml,
			expectedValues:         valuesYaml,
			errorMatcher:           nil,
			setup: func(chartDir string) error {
				files := []File{
					{
						path: filepath.Join(chartDir, "Chart.yaml"),
						data: chartYaml,
					},
					{
						path: filepath.Join(chartDir, "values.yaml"),
						data: valuesYaml,
					},
					{
						path: filepath.Join(chartDir, "templates/deployment.yaml"),
						data: deploymentYaml,
					},
				}

				return createFiles(files)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			// setup chartDir and dst as temporary directories.
			var chartDir, dst, filename string
			{
				chartDir, err = ioutil.TempDir(os.TempDir(), "architect-package-helmcharttask-chartdir-")
				if err != nil {
					t.Fatalf("error during chartDir directory creation: %v", err)
				}
				defer os.RemoveAll(chartDir)

				dst, err = ioutil.TempDir(os.TempDir(), "architect-package-helmcharttask-dst-")
				if err != nil {
					t.Fatalf("error during dst directory creation: %v", err)
				}
				defer os.RemoveAll(dst)
			}

			// setup Chart directories and files.
			if err = tc.setup(chartDir); err != nil {
				t.Fatalf("error during setup: %v", err)
			}

			// run the test: package the chart.
			task := NewPackageHelmChartTask(chartDir, dst)
			err = task.Run()

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if tc.expectedChartName == "" || tc.expectedChartVersion == "" {
				return
			}

			// packaged chart filename is deterministic.
			filename = filepath.Join(dst, fmt.Sprintf("%s-%s.tgz", tc.expectedChartName, tc.expectedChartVersion))

			// test for expected chart archive filename.
			_, err = os.Stat(filename)
			if err != nil {
				t.Fatalf("chart file (%s) not found: %v", filename, err)
			}

			// load the chart from the archive.
			chart, err := chartutil.LoadFile(filename)
			if err != nil {
				t.Fatalf("failed to load chart (%s): %v", filename, err)
			}
			metadata := chart.GetMetadata()
			values := chart.GetValues()

			// check for chart name.
			if name := metadata.GetName(); name != tc.expectedChartName {
				t.Fatalf("wrong chart name: expected %#q, got %#q", tc.expectedChartName, name)
			}

			// check for chart version.
			if version := metadata.GetVersion(); version != tc.expectedChartVersion {
				t.Fatalf("wrong chart version: expected %#q, got %#q", tc.expectedChartVersion, version)
			}

			if tc.expectedValues != "" {
				// check for chart values.
				if values := values.GetRaw(); values != tc.expectedValues {
					t.Fatalf("wrong chart values:\n>>> expected\n%#q\n>>> got\n%#q\n", tc.expectedValues, values)
				}
			}

			// check for chart deployment template.
			if tc.expectedDeploymentName != "" && tc.expectedDeploymentData != "" {
				var deploymentTemplate *hapichart.Template
				templates := chart.GetTemplates()
				for _, template := range templates {
					if name := template.GetName(); name == tc.expectedDeploymentName {
						deploymentTemplate = template
						break
					}
				}
				if deploymentTemplate == nil {
					t.Fatalf("not found deployment template: %#q", tc.expectedDeploymentName)
				}
				if data := deploymentTemplate.GetData(); !bytes.Equal(data, []byte(tc.expectedDeploymentData)) {
					t.Logf("%s", deploymentTemplate.GetData())
					t.Fatalf("wrong deployment template data:\n>>> expected\n%s\n>>> got\n%s\n", tc.expectedDeploymentData, data)
				}
			}
		})
	}
}

func createFiles(files []File) error {
	for _, file := range files {
		dir := filepath.Dir(file.path)
		if dir != "." && dir != "/" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return microerror.Mask(err)
			}
		}
		if err := ioutil.WriteFile(file.path, []byte(file.data), 0644); err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
