package main

import (
	"io"

	"github.com/spf13/cobra"
)

const controlDesc = `Check if artifact is allowed to be deployed.`

func newControlCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "control",
		Short: controlDesc,
		Long:  controlDesc,
	}

	// Add subcommands
	cmd.AddCommand(
		newControlDeploymentCmd(out),
		newControlPullRequestCmd(out),
	)

	return cmd
}