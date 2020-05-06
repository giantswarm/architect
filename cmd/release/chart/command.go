package chart

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/hook"
)

var (
	Cmd = &cobra.Command{
		Use:     "chart",
		Short:   "release chart as github release",
		RunE:    runReleaseError,
		PreRunE: hook.PreRunE,
	}
)
