package release

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/hook"
)

var (
	Cmd = &cobra.Command{
		Use:     "release",
		Short:   "release chart as github release",
		RunE:    runReleaseError,
		PreRunE: hook.PreRunE,
	}
)
