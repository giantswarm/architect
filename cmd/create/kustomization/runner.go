package kustomization

import (
	"fmt"
	"os"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type kustomization struct {
	Generators []string `json:"generators,omitempty"`
}

func runKustomizationError(cmd *cobra.Command, args []string) error {
	err := validateFlags()
	if err != nil {
		return microerror.Mask(err)
	}

	dirEntries, err := os.ReadDir(flag.Dir)
	if err != nil {
		return microerror.Mask(err)
	}

	k := kustomization{
		Generators: []string{},
	}

	for _, entry := range dirEntries {
		if !entry.Type().IsRegular() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		k.Generators = append(k.Generators, entry.Name())
	}

	data, err := yaml.Marshal(&k)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s", data)
	return nil
}
