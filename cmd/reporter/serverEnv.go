package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/merkely-development/reporter/internal/requests"
	"github.com/merkely-development/reporter/internal/server"
	"github.com/spf13/cobra"
)

const serverEnvDesc = `
List the artifacts deployed in a server environment and their digests 
and report them to Merkely. 
`

const serverEnvExample = `
* report directory artifacts running in a server at a list of paths:
merkely report env server prod --api-token 1234 --owner exampleOrg --id prod-server --paths a/b/c, e/f/g
`

type serverEnvOptions struct {
	paths   []string
	id      string
	verbose bool
}

func newServerEnvCmd(out io.Writer) *cobra.Command {
	o := new(serverEnvOptions)
	cmd := &cobra.Command{
		Use:     "server [-p /path/of/artifacts/directory] [-i infrastructure-identifier] env-name",
		Short:   "Report directory artifacts data in the given list of paths to Merkely.",
		Long:    serverEnvDesc,
		Aliases: []string{"directories"},
		Example: serverEnvExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("only environment name argument is allowed")
			}
			if len(args) == 0 || args[0] == "" {
				return fmt.Errorf("environment name is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			if o.id == "" {
				o.id = envName
			}

			url := fmt.Sprintf("%s/api/v1/environments/%s/%s/data", global.host, global.owner, envName)

			artifacts, err := server.CreateServerArtifactsData(o.paths, o.verbose)
			if err != nil {
				return err
			}
			requestBody := &requests.ServerEnvRequest{
				Artifacts: artifacts,
				Type:      "server",
				Id:        o.id,
			}
			js, _ := json.MarshalIndent(requestBody, "", "    ")

			return requests.SendPayload(js, url, global.apiToken,
				global.maxAPIRetries, global.dryRun)
		},
	}

	cmd.Flags().StringSliceVarP(&o.paths, "paths", "p", []string{}, "the comma separated list of artifact directories.")
	cmd.Flags().StringVarP(&o.id, "id", "i", "", "the unique identifier of the source infrastructure of the report (e.g. the K8S cluster/namespace name). If not set, it is defaulted to environment name.")
	cmd.Flags().BoolVarP(&o.verbose, "verbose", "v", false, "print verbose output of directory digest calculation.")
	return cmd
}