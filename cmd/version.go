package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	pkgproject "github.com/giantswarm/architect/pkg/project"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version information",
		Run:   runVersion,
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	if pkgproject.GitSHA() == "" && pkgproject.BuildTimestamp() == "" {
		fmt.Printf("version information not compiled\n")
		os.Exit(0)
	}

	fmt.Printf("Git Commit Hash: %s\n", pkgproject.GitSHA())
	fmt.Printf("Build Timestamp: %s\n", pkgproject.BuildTimestamp())
}
