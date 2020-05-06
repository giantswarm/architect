package release

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/hook"
)

var (
	Cmd = &cobra.Command{
		Use:     "release",
		Short:   "release operator versions or charts",
		RunE:    runReleaseError,
		PreRunE: hook.PreRunE,
	}
)
