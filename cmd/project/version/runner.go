package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runVersionError(cmd *cobra.Command, args []string) error {
	version := cmd.Flag("version").Value.String()

	fmt.Printf("%s\n", version)

	return nil
}
