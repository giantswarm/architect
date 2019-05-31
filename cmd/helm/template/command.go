package template

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "template",
		Short: "templates helm chart",
		RunE:  runTemplateError,
	}
)
