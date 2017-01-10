package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "architect",
	Short: "architect is a tool for managing Giant Swarm release engineering",
}
