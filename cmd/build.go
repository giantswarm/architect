package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/build"
)

func init() {
	RootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Test and build the project",
	Run:   runBuildCmd,
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	log.Printf("Configuring build\n")
	buildConfig := build.Config{
		Type: build.GolangType,
	}

	log.Printf("Constructing builder\n")
	builder := build.New(buildConfig)

	log.Printf("Running tests\n")
	if err := builder.Test(); err != nil {
		log.Fatalf("An error occurred during test: %v\n", err)
	}

	log.Printf("Building project\n")
	if err := builder.Build(); err != nil {
		log.Fatalf("An error occurred during build: %v\n", err)
	}
}
