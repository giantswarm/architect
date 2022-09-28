package kustomization

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_runKustomizationError(t *testing.T) {
	testCases := []struct {
		name      string
		inputPath string
		expected  string
	}{
		{
			name:      "flawless",
			inputPath: "testdata/input/flux-manifests_0",
			expected:  "testdata/expected/kustomization.yaml_0",
		},
		{
			name:      "flawless with existing kustomization.yaml",
			inputPath: "testdata/input/flux-manifests_1",
			expected:  "testdata/expected/kustomization.yaml_0",
		},
		{
			name:      "flawless with patches",
			inputPath: "testdata/input/flux-manifests_2",
			expected:  "testdata/expected/kustomization.yaml_1",
		},
		{
			name:      "empty dir",
			inputPath: "testdata/input/flux-manifests_3",
			expected:  "testdata/expected/kustomization.yaml_2",
		},
		{
			name:      "existing and complicated kustomization.yaml",
			inputPath: "testdata/input/flux-manifests_4",
			expected:  "testdata/expected/kustomization.yaml_3",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			var err error

			out := new(bytes.Buffer)

			cmd := NewCommand()
			cmd.SetOut(out)

			flag.Dir = tc.inputPath

			err = runKustomizationError(cmd, []string{})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			exp, err := os.ReadFile(tc.expected)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			diff := cmp.Diff(string(exp), out.String())
			if diff != "" {
				t.Fatalf("files do not match, got:\n %s", diff)
			}
		})
	}
}
