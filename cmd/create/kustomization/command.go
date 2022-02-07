package kustomization

import (
	"fmt"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

var flag struct {
	Dir string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kustomization",
		Short: "Create Kustomization for all Flux Generators",
		RunE:  runKustomizationError,
	}

	cmd.Flags().StringVar(&flag.Dir, "dir", "", "App collection directory.")

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
