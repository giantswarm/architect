package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/workflow"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build the specified project",
		Run:   runBuildCmd,
	}

	dryRun bool
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVar(&dryRun, "dry-run", false, "do not perform any actions")
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	executor := workflow.GetExecutor(dryRun)

	build := workflow.NewGolangWorkflow()
	if err := workflow.RunWorkflow(executor, build); err != nil {
		log.Fatalln(err)
	}
}
