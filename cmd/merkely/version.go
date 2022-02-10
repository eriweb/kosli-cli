package main

import (
	"fmt"
	"io"

	"github.com/merkely-development/reporter/internal/version"
	"github.com/spf13/cobra"
)

const versionDesc = `
Print the version for Merkely CLI.

The output will look something like this:
version.BuildInfo{Version:"v0.0.1", GitCommit:"fe51cd1e31e6a202cba7dead9552a6d418ded79a", GitTreeState:"clean", GoVersion:"go1.16.3"}

- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- GitTreeState is "clean" if there are no local code changes when this binary was
  built, and "dirty" if the binary was built from locally modified code.
- GoVersion is the version of Go that was used to compile Merkely CLI.
`

type versionOptions struct {
	short bool
}

func newVersionCmd(out io.Writer) *cobra.Command {
	o := new(versionOptions)
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client version information",
		Long:  versionDesc,
		Args:  NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}

	cmd.Flags().BoolVarP(&o.short, "short", "s", false, "only print the version number")

	return cmd
}

func (o *versionOptions) run(out io.Writer) error {
	fmt.Fprint(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	if short {
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", version.Get())
}
