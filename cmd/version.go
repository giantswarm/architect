package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version information",
		Run:   runVersion,
	}

	Commit         string
	BuildTimestamp string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	if Commit == "" && BuildTimestamp == "" {
		fmt.Printf("version information not compiled\n")
		os.Exit(0)
	}

	fmt.Printf("Git Commit Hash: %s\n", Commit)
	fmt.Printf("Build Timestamp: %s\n", BuildTimestamp)
}
