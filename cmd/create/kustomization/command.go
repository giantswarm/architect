package kustomization

import (
	"fmt"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

var flag struct {
	Dir        string
	Generators bool
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kustomization",
		Short: "Create Kustomization for all files in a directory",
		RunE:  runKustomizationError,
	}

	cmd.Flags().StringVar(&flag.Dir, "dir", "", "App collection directory.")
	cmd.Flags().BoolVar(&flag.Generators, "generators", true, "When enabled, the directory content "+
		"is assumed to be generators and will be listed under .generators in the resulted kustomization.yaml."+
		"Will render the list under .resources otherwise. Default is true for backward compatibility.")

	return cmd
}

func validateFlags() error {
	var errors []string

	if flag.Dir == "" {
		errors = append(errors, "--dir is required")
	}
	if len(errors) != 0 {
		return microerror.Mask(fmt.Errorf("invalid flag(s): %s", strings.Join(errors, ", ")))
	}

	return nil
}
