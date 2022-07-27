package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kosli-dev/cli/internal/requests"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const artifactGetDesc = `Get artifact from specified pipeline`

type artifactGetOptions struct {
	json         bool
	pipelineName string
}

func newArtifactGetCmd(out io.Writer) *cobra.Command {
	o := new(artifactGetOptions)
	cmd := &cobra.Command{
		Use:   "get ARTIFACT-DIGEST",
		Short: artifactGetDesc,
		Long:  artifactGetDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Owner", "ApiToken"})
			if err != nil {
				return ErrorAfterPrintingHelp(cmd, err.Error())
			}
			if len(args) < 1 {
				return ErrorAfterPrintingHelp(cmd, "pipeline name argument is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out, args)
		},
	}

	cmd.Flags().StringVarP(&o.pipelineName, "pipeline", "p", "", pipelineNameFlag)
	cmd.Flags().BoolVarP(&o.json, "json", "j", false, jsonOutputFlag)

	err := RequireFlags(cmd, []string{"pipeline"})
	if err != nil {
		log.Fatalf("failed to configure required flags: %v", err)
	}

	return cmd
}

func (o *artifactGetOptions) run(out io.Writer, args []string) error {
	url := fmt.Sprintf("%s/api/v1/projects/%s/%s/artifacts/%s", global.Host, global.Owner, o.pipelineName, args[0])
	response, err := requests.DoBasicAuthRequest([]byte{}, url, "", global.ApiToken,
		global.MaxAPIRetries, http.MethodGet, map[string]string{}, logrus.New())

	if err != nil {
		return err
	}

	if o.json {
		pj, err := prettyJson(response.Body)
		if err != nil {
			return err
		}
		fmt.Println(pj)
		return nil
	}

	var artifact map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &artifact)
	if err != nil {
		return err
	}

	artifactData := artifact["evidence"].(map[string]interface{})["artifact"].(map[string]interface{})
	fmt.Printf("Name: %s\n", artifactData["filename"].(string))
	fmt.Printf("State: %s\n", artifact["state"].(string))
	fmt.Printf("Git commit: %s\n", artifactData["git_commit"].(string))
	fmt.Printf("Build URL: %s\n", artifactData["build_url"].(string))
	fmt.Printf("Commit URL: %s\n", artifactData["commit_url"].(string))
	fmt.Printf("Created at: %f\n", artifactData["logged_at"].(float64))
	// fmt.Printf("Description: %s\n", pipeline["description"])
	// fmt.Printf("Visibility: %s\n", pipeline["visibility"])
	// template := fmt.Sprintf("%s", pipeline["template"])
	// template = strings.Replace(template, " ", ", ", -1)
	// fmt.Printf("Template: %s\n", template)
	// timeago.English.Max = 36 * timeago.Month
	// timestampFloat, err := strconv.ParseFloat(pipeline["last_deployment_at"].(string), 64)
	// if err != nil {
	// 	return err
	// }
	// timestamp := time.Unix(int64(timestampFloat), 0)
	// last_deployment_at := timeago.English.Format(timestamp)
	// fmt.Printf("Last deployment at: %s \u2022 %s\n", timestamp.Format(time.RFC822), last_deployment_at)

	return nil
}
