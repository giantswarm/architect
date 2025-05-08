package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	pkgproject "github.com/giantswarm/architect/v2/pkg/project"
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

	fmt.Printf(
		"Version: %s\nGit Commit Hash: %s\nBuild Timestamp: %s\n",
		pkgproject.Version(),
		pkgproject.GitSHA(),
		pkgproject.BuildTimestamp(),
	)
}
