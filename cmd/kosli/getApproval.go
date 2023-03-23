package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kosli-dev/cli/internal/output"
	"github.com/kosli-dev/cli/internal/requests"
	"github.com/spf13/cobra"
)

const getApprovalShortDesc = `Get an approval from a specified flow.`
const getApprovalLongDesc = getApprovalShortDesc + `
The expected argument is an expression to specify the approval to get.
It has the format <FLOW_NAME>[SEPARATOR][INTEGER_REFERENCE]

Specify SNAPPISH by:
	flowName~<N>  N'th behind the latest approval
	flowName#<N>  approval number N
	flowName      the latest approval

Examples of valid expressions are: flow (latest approval), flow#10 (approval number 10), flow~2 (the third latest approval)`

const getApprovalExample = `
# get second behind the latest approval from a flow
kosli get approval flowName~1 \
	--api-token yourAPIToken \
	--org orgName

# get the 10th approval from a flow
kosli get approval flowName#10 \
	--api-token yourAPIToken \
	--org orgName

# get the latest approval from a flow
kosli get approval flowName \
	--api-token yourAPIToken \
	--org orgName`

type getApprovalOptions struct {
	output string
}

func newGetApprovalCmd(out io.Writer) *cobra.Command {
	o := new(getApprovalOptions)
	cmd := &cobra.Command{
		Use:     "approval SNAPPISH",
		Short:   getApprovalShortDesc,
		Long:    getApprovalLongDesc,
		Example: getApprovalExample,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Org", "ApiToken"})
			if err != nil {
				return ErrorBeforePrintingUsage(cmd, err.Error())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out, args)
		},
	}

	cmd.Flags().StringVarP(&o.output, "output", "o", "table", outputFlag)
	return cmd
}

func (o *getApprovalOptions) run(out io.Writer, args []string) error {
	url := fmt.Sprintf("%s/api/v1/projects/%s/approval/?snappish=%s", global.Host, global.Org, url.QueryEscape(args[0]))

	reqParams := &requests.RequestParams{
		Method:   http.MethodGet,
		URL:      url,
		Password: global.ApiToken,
	}
	response, err := kosliClient.Do(reqParams)
	if err != nil {
		return err
	}

	return output.FormattedPrint(response.Body, o.output, out, 0,
		map[string]output.FormatOutputFunc{
			"table": printApprovalAsTable,
			"json":  output.PrintJson,
		})
}

func printApprovalAsTable(raw string, out io.Writer, page int) error {
	var approval map[string]interface{}
	err := json.Unmarshal([]byte(raw), &approval)
	if err != nil {
		return err
	}

	rows := []string{}
	rows = append(rows, fmt.Sprintf("ID:\t%d", int64(approval["release_number"].(float64))))
	rows = append(rows, fmt.Sprintf("Artifact fingerprint:\t%s", approval["base_artifact"].(string)))
	rows = append(rows, fmt.Sprintf("Artifact name:\t%s", approval["artifact_name"].(string)))
	rows = append(rows, fmt.Sprintf("State:\t%s", approval["state"].(string)))
	lastModifiedAt, err := formattedTimestamp(approval["last_modified_at"], false)
	if err != nil {
		return err
	}
	rows = append(rows, fmt.Sprintf("Last modified at:\t%s", lastModifiedAt))
	reviews := approval["approvals"].([]interface{})
	if len(reviews) > 0 {
		rows = append(rows, "Reviews:")
		for _, review := range reviews {
			convertedReview := review.(map[string]interface{})
			approvedBy := "Unknown"
			if convertedReview["approved_by"] != nil {
				approvedBy = convertedReview["approved_by"].(string)
			}
			createdAt, err := formattedTimestamp(convertedReview["timestamp"], true)
			if err != nil {
				return err
			}
			reviewRow := fmt.Sprintf("\t%s By: %s on %s", convertedReview["state"].(string), approvedBy, createdAt)
			rows = append(rows, reviewRow)
		}
	} else {
		rows = append(rows, "Reviews:\tNone")
	}

	commits := approval["src_commit_list"].([]interface{})
	if len(reviews) > 0 {
		rows = append(rows, "Changes:")
		for _, commit := range commits {
			convertedCommit := commit.(map[string]interface{})
			commitRow := fmt.Sprintf("\tGit commit:%s", convertedCommit["commit_sha"].(string))
			rows = append(rows, commitRow)
			artifact_digests := convertedCommit["artifact_digests"].([]interface{})
			if len(artifact_digests) == 0 {
				commitRow = "\tNo artifacts produced from this commit"
				rows = append(rows, commitRow)
			} else {
				commitRow = "\tProduced artifact fingerprint(s):"
				rows = append(rows, commitRow)
				for _, digest := range artifact_digests {
					digestRow := fmt.Sprintf("\t\t%s", digest)
					rows = append(rows, digestRow)
				}

			}
		}
	} else {
		rows = append(rows, "Changes:\tNone")
	}

	tabFormattedPrint(out, []string{}, rows)
	return nil
}