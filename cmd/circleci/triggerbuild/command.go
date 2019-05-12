package triggerbuild

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "trigger-build",
		Short: "Trigger a CircleCI build",
		RunE:  runTriggerBuildError,
	}
)
