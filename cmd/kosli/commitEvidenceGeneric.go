package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kosli-dev/cli/internal/requests"
	"github.com/spf13/cobra"
)

type CommitEvidenceGenericPayload struct {
	CommitSHA           string      `json:"commit_sha"`
	Pipelines           []string    `json:"pipelines,omitempty"`
	Description         string      `json:"description,omitempty"`
	Compliant           bool        `json:"is_compliant"`
	EvidenceName        string      `json:"name"`
	BuildUrl            string      `json:"build_url"`
	EvidenceUrl         string      `json:"evidence_url,omitempty"`
	EvidenceFingerprint string      `json:"evidence_fingerprint,omitempty"`
	UserData            interface{} `json:"user_data,omitempty"`
}

type genericCommitEvidenceOptions struct {
	userDataFile string
	payload      CommitEvidenceGenericPayload
}

const genericCommitEvidenceShortDesc = `Report Generic evidence for a commit in a Kosli pipeline.`

const genericCommitEvidenceLongDesc = genericCommitEvidenceShortDesc

const genericCommitEvidenceExample = `
# report Generic evidence for a commit related to one Kosli pipeline:
kosli commit report evidence generic \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--description "some description" \
	--compliant \
	--pipelines yourPipelineName \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--owner yourOrgName

# report Generic evidence for a commit related to multiple Kosli pipelines with user-data:
kosli commit report evidence generic \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--description "some description" \
	--compliant \
	--pipelines yourFirstPipelineName,yourSecondPipelineName \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--owner yourOrgName \
	--user-data /path/to/json/file.json
`

func newGenericCommitEvidenceCmd(out io.Writer) *cobra.Command {
	o := new(genericCommitEvidenceOptions)
	cmd := &cobra.Command{
		Use:     "generic",
		Short:   genericCommitEvidenceShortDesc,
		Long:    genericCommitEvidenceLongDesc,
		Example: genericCommitEvidenceExample,
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Owner", "ApiToken"})
			if err != nil {
				return ErrorBeforePrintingUsage(cmd, err.Error())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(args)
		},
	}

	ci := WhichCI()
	cmd.Flags().StringVar(&o.payload.CommitSHA, "commit", DefaultValue(ci, "git-commit"), evidenceCommitFlag)
	cmd.Flags().StringSliceVarP(&o.payload.Pipelines, "pipelines", "p", []string{}, pipelinesFlag)
	cmd.Flags().StringVarP(&o.payload.BuildUrl, "build-url", "b", DefaultValue(ci, "build-url"), evidenceBuildUrlFlag)
	cmd.Flags().BoolVarP(&o.payload.Compliant, "compliant", "C", false, evidenceCompliantFlag)
	cmd.Flags().StringVarP(&o.payload.Description, "description", "d", "", evidenceDescriptionFlag)
	cmd.Flags().StringVarP(&o.payload.EvidenceName, "name", "n", "", evidenceNameFlag)
	cmd.Flags().StringVar(&o.payload.EvidenceUrl, "evidence-url", "", evidenceUrlFlag)
	cmd.Flags().StringVar(&o.payload.EvidenceFingerprint, "evidence-fingerprint", "", evidenceFingerprintFlag)
	cmd.Flags().StringVarP(&o.userDataFile, "user-data", "u", "", evidenceUserDataFlag)
	addDryRunFlag(cmd)

	err := RequireFlags(cmd, []string{"commit", "build-url", "name"})
	if err != nil {
		logger.Error("failed to configure required flags: %v", err)
	}

	return cmd
}

func (o *genericCommitEvidenceOptions) run(args []string) error {
	var err error
	url := fmt.Sprintf("%s/api/v1/projects/%s/commit/evidence/generic", global.Host, global.Owner)
	o.payload.UserData, err = LoadJsonData(o.userDataFile)
	if err != nil {
		return err
	}

	reqParams := &requests.RequestParams{
		Method:   http.MethodPut,
		URL:      url,
		Payload:  o.payload,
		DryRun:   global.DryRun,
		Password: global.ApiToken,
	}
	_, err = kosliClient.Do(reqParams)
	if err == nil && !global.DryRun {
		logger.Info("generic evidence '%s' is reported to commit: %s", o.payload.EvidenceName, o.payload.CommitSHA)
	}
	return err
}
