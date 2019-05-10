package circleci

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "circleci",
		Short: "interacts with CircleCI",
	}
)
