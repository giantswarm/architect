package pipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/architect/pipeline"
	"github.com/spf13/afero"
)

func TestStartChannel(t *testing.T) {
	fs := afero.NewMemMapFs()
	project := "myproject"
	workingDirectory := "/myworkingDirectory"

	chartDirectory := filepath.Join(workingDirectory, "helm", project+"-chart")
	chartFile := filepath.Join(chartDirectory, "Chart.yaml")

	err := fs.MkdirAll(chartDirectory, os.ModePerm)
	if err != nil {
		t.Fatalf("could not create %q", workingDirectory)
	}

	tcs := []struct {
		description      string
		chartFileContent string
		fileDoesntExist  bool
		expectedError    bool
		expectedChannel  string
	}{
		{
			description:     "chart file doesn't exists",
			fileDoesntExist: true,
		},
		{
			description: "yaml error",
			chartFileContent: `version: 1.2.3
  base-field: value1
    well-indented-field: value2
   bad-indented-field: value3
   `,
			expectedError: true,
		},
		{
			description:      "base case",
			chartFileContent: "version: 1.2.3",
			expectedChannel:  "1-2-beta",
		},
		{
			description:      "do not include initial v",
			chartFileContent: "version: v1.2.3",
			expectedChannel:  "1-2-beta",
		},
		{
			description:      "manage non semantic versions",
			chartFileContent: "version: v1.v2wer.abc",
			expectedChannel:  "1-v2wer-beta",
		},
		{
			description:      "fail if there are no identifiable major and minor",
			chartFileContent: "version: v1v2wer",
			expectedError:    true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			if tc.fileDoesntExist {
				_, err := pipeline.StartChannel(fs, workingDirectory, project)
				if err == nil {
					t.Fatalf("expected error didn't happen")
				}
				return
			}

			err := afero.WriteFile(fs, chartFile, []byte(tc.chartFileContent), os.ModePerm)
			if err != nil {
				t.Errorf("could not create file %q", chartFile)
			}

			channel, err := pipeline.StartChannel(fs, workingDirectory, project)
			if tc.expectedError && err == nil {
				t.Errorf("expected error didn't happen")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("could not get start channel %v", err)
			}

			if channel != tc.expectedChannel {
				t.Errorf("didn't get expected channel, want %q, got %q", tc.expectedChannel, channel)
			}
		})
	}
}

func TestEndChannel(t *testing.T) {
	fs := afero.NewMemMapFs()
	project := "myproject"
	workingDirectory := "/myworkingDirectory"

	chartDirectory := filepath.Join(workingDirectory, "helm", project+"-chart")
	chartFile := filepath.Join(chartDirectory, "Chart.yaml")

	err := fs.MkdirAll(chartDirectory, os.ModePerm)
	if err != nil {
		t.Fatalf("could not create %q", workingDirectory)
	}

	tcs := []struct {
		description      string
		chartFileContent string
		fileDoesntExist  bool
		expectedError    bool
		expectedChannel  string
	}{
		{
			description:     "chart file doesn't exists",
			fileDoesntExist: true,
		},
		{
			description: "yaml error",
			chartFileContent: `version: 1.2.3
  base-field: value1
    well-indented-field: value2
   bad-indented-field: value3
   `,
			expectedError: true,
		},
		{
			description:      "base case",
			chartFileContent: "version: 1.2.3",
			expectedChannel:  "1-2-stable",
		},
		{
			description:      "do not include initial v",
			chartFileContent: "version: v1.2.3",
			expectedChannel:  "1-2-stable",
		},
		{
			description:      "manage non semantic versions",
			chartFileContent: "version: v1.v2wer.abc",
			expectedChannel:  "1-v2wer-stable",
		},
		{
			description:      "fail if there are no identifiable major and minor",
			chartFileContent: "version: v1v2wer",
			expectedError:    true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			if tc.fileDoesntExist {
				_, err := pipeline.StartChannel(fs, workingDirectory, project)
				if err == nil {
					t.Fatalf("expected error didn't happen")
				}
				return
			}

			err := afero.WriteFile(fs, chartFile, []byte(tc.chartFileContent), os.ModePerm)
			if err != nil {
				t.Errorf("could not create file %q", chartFile)
			}

			channel, err := pipeline.EndChannel(fs, workingDirectory, project)
			if tc.expectedError && err == nil {
				t.Errorf("expected error didn't happen")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("could not get start channel %v", err)
			}

			if channel != tc.expectedChannel {
				t.Errorf("didn't get expected channel, want %q, got %q", tc.expectedChannel, channel)
			}
		})
	}
}
