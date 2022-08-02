package main

import (
	"io"

	"github.com/spf13/cobra"
)

const artifactDesc = `All artifacts report operations in a Kosli pipeline.`

func newArtifactCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: artifactDesc,
		Long:  artifactDesc,
	}

	// Add subcommands
	cmd.AddCommand(
		newArtifactReportCmd(out),
		newArtifactLsCmd(out),
	)

	return cmd
}