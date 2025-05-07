package main

import (
	"fmt"
	"os"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/v2/cmd"
)

func main() {
	cmd.RootCmd.SilenceErrors = true
	cmd.RootCmd.SilenceUsage = true
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", microerror.Pretty(err, true))
		os.Exit(2)
	}
}
