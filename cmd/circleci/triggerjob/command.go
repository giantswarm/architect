package triggerjob

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "trigger-job",
		Short: "Trigger a CircleCI job",
		RunE:  runTriggerJobError,
	}
)
