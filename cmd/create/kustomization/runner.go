package kustomization

import (
	"fmt"
	"os"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func runKustomizationError(cmd *cobra.Command, args []string) error {
	err := validateFlags()
	if err != nil {
		return microerror.Mask(err)
	}

	dirEntries, err := os.ReadDir(flag.Dir)
	if err != nil {
		return microerror.Mask(err)
	}

	kusRaw, err := os.ReadFile(fmt.Sprintf("%s/kustomization.yaml", flag.Dir))
	if err != nil && !os.IsNotExist(err) {
		return microerror.Mask(err)
	}

	kus := make(map[string]interface{})

	err = yaml.Unmarshal(kusRaw, &kus)
	if err != nil {
		return microerror.Mask(err)
	}

	files := make([]string, 0)
	for _, entry := range dirEntries {
		if !entry.Type().IsRegular() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		if entry.Name() == "kustomization.yaml" {
			continue
		}
		files = append(files, entry.Name())
	}

	if flag.Generators {
		kus["generators"] = files
	} else {
		kus["resources"] = files
	}

	data, err := yaml.Marshal(&kus)
	if err != nil {
		return microerror.Mask(err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s", data)
	return nil
}
