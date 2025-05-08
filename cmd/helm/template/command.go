package template

import (
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/v2/cmd/hook"
)

var (
	Cmd = &cobra.Command{
		Use:     "template",
		Short:   "templates helm chart",
		RunE:    runTemplateError,
		PreRunE: hook.PreRunE,
	}
)
