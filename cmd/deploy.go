package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the project",
	Run:   runDeployCmd,
}

func runDeployCmd(cmd *cobra.Command, args []string) {

}
