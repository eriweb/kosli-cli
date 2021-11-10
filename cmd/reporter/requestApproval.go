package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/merkely-development/reporter/internal/digest"
	"github.com/merkely-development/reporter/internal/requests"
	"github.com/spf13/cobra"
)

func newRequestApprovalCmd(out io.Writer) *cobra.Command {
	o := new(approvalOptions)
	cmd := &cobra.Command{
		Use:   "request ARTIFACT-NAME-OR-PATH",
		Short: "Request an approval for deploying an artifact in Merkely. ",
		Long:  requestApprovalDesc(),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("only one argument (docker image name or file/dir path) is allowed")
			}
			if len(args) == 0 || args[0] == "" {
				return fmt.Errorf("docker image name or file/dir path is required")
			}

			if o.artifactType == "" && o.inputSha256 == "" {
				return fmt.Errorf("either --type or --sha256 must be specified")
			}

			err := RequireGlobalFlags(global, []string{"Owner", "ApiToken"})
			if err != nil {
				return err
			}

			if o.inputSha256 != "" {
				if err := digest.ValidateDigest(o.inputSha256); err != nil {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if o.inputSha256 != "" {
				o.payload.ArtifactSha256 = o.inputSha256
			} else {
				o.payload.ArtifactSha256, err = GetSha256Digest(o.artifactType, args[0])
				if err != nil {
					return err
				}
			}

			url := fmt.Sprintf("%s/api/v1/projects/%s/%s/approvals/", global.Host, global.Owner, o.pipelineName)
			o.payload.Reviews = []map[string]string{}
			o.payload.UserData, err = LoadUserData(o.userDataFile)
			if err != nil {
				return err
			}
			o.payload.CommitList, err = listCommitsBetween(o.srcRepoRoot, o.oldestSrcCommit, o.newestSrcCommit)
			if err != nil {
				return err
			}
			js, _ := json.MarshalIndent(o.payload, "", "    ")

			return requests.SendPayload(js, url, global.ApiToken,
				global.MaxAPIRetries, global.DryRun, "POST", log)
		},
	}

	cmd.Flags().StringVarP(&o.artifactType, "artifact-type", "t", "", "The type of the artifact to be approved. Options are [dir, file, docker]. Only required if you don't specify --sha256.")
	cmd.Flags().StringVarP(&o.inputSha256, "sha256", "s", "", "The SHA256 fingerprint for the artifact to be approved. Only required if you don't specify --type.")
	cmd.Flags().StringVarP(&o.pipelineName, "pipeline", "p", "", "The Merkely pipeline name.")
	cmd.Flags().StringVarP(&o.payload.Description, "description", "d", "", "[optional] The approval description.")
	cmd.Flags().StringVarP(&o.userDataFile, "user-data", "u", "", "[optional] The path to a JSON file containing additional data you would like to attach to this approval.")
	cmd.Flags().StringVar(&o.oldestSrcCommit, "oldest-commit", "", "The source commit sha for the oldest change in the deployment approval.")
	cmd.Flags().StringVar(&o.newestSrcCommit, "newest-commit", "HEAD", "The source commit sha for the newest change in the deployment approval.")
	cmd.Flags().StringVar(&o.srcRepoRoot, "repo-root", "/src", "The directory where the source git repository is volume-mounted.")

	err := RequireFlags(cmd, []string{"pipeline", "oldest-commit"})
	if err != nil {
		log.Fatalf("failed to configure required flags: %v", err)
	}

	return cmd
}

func requestApprovalDesc() string {
	return `
   Request an approval of a deployment of an artifact in Merkely. The request should be reviewed in Merkely UI.
   The artifact SHA256 fingerprint is calculated or alternatively it can be provided directly. 
   `
}
