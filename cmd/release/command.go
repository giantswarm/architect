package release

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "release",
		Short: "release chart as github release",
		RunE:  runReleaseError,
	}
)
