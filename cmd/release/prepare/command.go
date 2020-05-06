package prepare

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/hook"
)

var (
	Cmd = &cobra.Command{
		Use:     "prepare",
		Short:   "prepare changelog and operator version to be released",
		RunE:    runReleaseError,
		PreRunE: hook.PreRunE,
	}
)
