package version

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/hook"
)

var (
	Cmd = &cobra.Command{
		Use:     "version",
		Short:   "show project version",
		RunE:    runVersionError,
		PreRunE: hook.PreRunE,
	}
)
