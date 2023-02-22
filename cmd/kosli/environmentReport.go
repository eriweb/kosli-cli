package main

import (
	"io"

	"github.com/spf13/cobra"
)

const environmentReportDesc = `Report artifacts running in an environment to Kosli. `

func newEnvironmentReportCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: environmentReportDesc,
		Long:  environmentReportDesc,
	}

	// Add subcommands
	cmd.AddCommand(
		newEnvironmentReportServerCmd(out),
		newEnvironmentReportS3Cmd(out),
		newEnvironmentReportLambdaCmd(out),
	)

	return cmd
}
