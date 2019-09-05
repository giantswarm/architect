package appcr

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "appcr",
		Short: "create App CR resource.",
		RunE:  runAppCRError,
	}
)
