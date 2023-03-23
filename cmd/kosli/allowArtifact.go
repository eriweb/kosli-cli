package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/kosli-dev/cli/internal/requests"
	"github.com/spf13/cobra"
)

type allowArtifactOptions struct {
	fingerprintOptions *fingerprintOptions
	payload            AllowlistPayload
}

type AllowlistPayload struct {
	Fingerprint string `json:"sha256"`
	Filename    string `json:"artifact_name"`
	Reason      string `json:"description"`
	Environment string `json:"environment_name"`
}

const allowArtifactShortDesc = `Add an artifact to an environment's allowlist.`

const allowArtifactLongDesc = allowArtifactShortDesc + `
` + fingerprintDesc

func newAllowArtifactCmd(out io.Writer) *cobra.Command {
	o := new(allowArtifactOptions)
	o.fingerprintOptions = new(fingerprintOptions)
	cmd := &cobra.Command{
		Use:   "artifact {IMAGE-NAME | FILE-PATH | DIR-PATH}",
		Short: allowArtifactShortDesc,
		Long:  allowArtifactLongDesc,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Org", "ApiToken"})
			if err != nil {
				return ErrorBeforePrintingUsage(cmd, err.Error())
			}

			err = ValidateArtifactArg(args, o.fingerprintOptions.artifactType, o.payload.Fingerprint, true)
			if err != nil {
				return ErrorBeforePrintingUsage(cmd, err.Error())
			}
			return ValidateRegistryFlags(cmd, o.fingerprintOptions)

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(args)
		},
	}

	cmd.Flags().StringVarP(&o.payload.Fingerprint, "fingerprint", "F", "", fingerprintFlag)
	cmd.Flags().StringVarP(&o.payload.Environment, "environment", "e", "", envAllowListFlag)
	cmd.Flags().StringVar(&o.payload.Reason, "reason", "", reasonFlag)
	addFingerprintFlags(cmd, o.fingerprintOptions)
	addDryRunFlag(cmd)

	err := RequireFlags(cmd, []string{"environment", "reason"})
	if err != nil {
		logger.Error("failed to configure required flags: %v", err)
	}

	return cmd
}

func (o *allowArtifactOptions) run(args []string) error {
	if o.payload.Fingerprint != "" {
		o.payload.Filename = args[0]
	} else {
		var err error
		o.payload.Fingerprint, err = GetSha256Digest(args[0], o.fingerprintOptions, logger)
		if err != nil {
			return err
		}
		if o.fingerprintOptions.artifactType == "dir" || o.fingerprintOptions.artifactType == "file" {
			o.payload.Filename = filepath.Base(args[0])
		} else {
			o.payload.Filename = args[0]
		}
	}

	url := fmt.Sprintf("%s/api/v1/policies/%s/allowedartifacts/", global.Host, global.Org)

	reqParams := &requests.RequestParams{
		Method:   http.MethodPut,
		URL:      url,
		Payload:  o.payload,
		DryRun:   global.DryRun,
		Password: global.ApiToken,
	}
	_, err := kosliClient.Do(reqParams)
	if err == nil && !global.DryRun {
		logger.Info("artifact %s was allow listed in environment: %s", o.payload.Fingerprint, o.payload.Environment)
	}
	return err
}