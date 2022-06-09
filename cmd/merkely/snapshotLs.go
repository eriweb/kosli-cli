package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/merkely-development/reporter/internal/requests"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const snapshotLsDesc = `
List snapshot.
`

type snapshotLsOptions struct {
	// long bool
}

type Annotation struct {
	Type string `json:"type"`
	Was  int
	Now  int
}

type Owner struct {
	ApiVersion         string
	Kind               string
	Name               string
	Uid                string
	Controller         bool
	BlockOwnerDeletion bool
}

type PodContent struct {
	Namespace         string
	CreationTimestamp int64
	Owners            []Owner
}

type Artifact struct {
	Name              string
	Pipeline_name     string
	Compliant         bool
	Deployments       []int
	Sha256            string
	CreationTimestamp []int64
	Pods              map[string]PodContent
	Annotation        Annotation
}

type Snapshot struct {
	Index     int
	Timestamp float32
	Type      string `json:"type"`
	User_id   string
	User_name string
	Artifacts []Artifact
	Compliant bool
}

type SnapshotType struct {
	Type string `json:"type"`
}

func newSnapshotLsCmd(out io.Writer) *cobra.Command {
	o := new(snapshotLsOptions)
	cmd := &cobra.Command{
		Use:   "snap",
		Short: snapshotLsDesc,
		Long:  snapshotLsDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out, args)
		},
	}

	// cmd.Flags().BoolVarP(&o.long, "long", "l", false, environmentLongFlag)

	return cmd
}

func (o *snapshotLsOptions) run(out io.Writer, args []string) error {
	url := fmt.Sprintf("%s/api/v1/environments/%s/%s/data", global.Host, global.Owner, args[0])
	response, err := requests.DoBasicAuthRequest([]byte{}, url, "", global.ApiToken,
		global.MaxAPIRetries, http.MethodGet, map[string]string{}, logrus.New())

	if err != nil {
		return fmt.Errorf("merkely server %s is unresponsive", global.Host)
	}

	var snapshotType SnapshotType
	err = json.Unmarshal([]byte(response.Body), &snapshotType)
	if err != nil {
		return err
	}

	if snapshotType.Type == "K8S" || snapshotType.Type == "ECS" {
		return showK8sEcs(response)
	} else if snapshotType.Type == "server" {
		return showServer(response)
	}
	return nil
}

func showK8sEcs(response *requests.HTTPResponse) error {
	var snapshot Snapshot
	err := json.Unmarshal([]byte(response.Body), &snapshot)
	if err != nil {
		return err
	}

	formatStringHead := "%-7s  %-40s  %-10s  %-17s  %-25s  %-10s\n"
	formatStringLine := "%-7s  %-40s  %-10s  %-17s  %-25s  %-10d\n"
	fmt.Printf(formatStringHead, "COMMIT", "IMAGE", "TAG", "SHA256", "SINCE", "REPLICAS")

	for _, artifact := range snapshot.Artifacts {
		since := time.Unix(artifact.CreationTimestamp[0], 0).Format(time.RFC3339)
		artifactNameSplit := strings.Split(artifact.Name, ":")
		artifactName := artifactNameSplit[0]
		if len(artifactName) > 40 {
			artifactName = artifactName[:18] + "..." + artifactName[len(artifactName)-19:]
		}
		artifactTag := ""
		if len(artifactNameSplit) > 1 {
			artifactTag = artifactNameSplit[1]
			if len(artifactTag) > 10 {
				artifactTag = artifactTag[:10]
			}
		}
		shortSha := ""
		if len(artifact.Sha256) == 64 {
			shortSha = artifact.Sha256[:7] + "..." + artifact.Sha256[64-7:]
		}
		fmt.Printf(formatStringLine, "xxxx", artifactName, artifactTag, shortSha, since, len(artifact.CreationTimestamp))
	}

	return nil
}

func showServer(response *requests.HTTPResponse) error {
	var snapshot Snapshot
	err := json.Unmarshal([]byte(response.Body), &snapshot)
	if err != nil {
		return err
	}

	formatStringHead := "%-7s  %-40s  %-17s  %-25s  %-10s\n"
	formatStringLine := "%-7s  %-40s  %-17s  %-25s  %-10d\n"
	fmt.Printf(formatStringHead, "COMMIT", "IMAGE", "SHA256", "SINCE", "REPLICAS")

	for _, artifact := range snapshot.Artifacts {
		since := time.Unix(artifact.CreationTimestamp[0], 0).Format(time.RFC3339)
		artifactName := artifact.Name
		if len(artifactName) > 40 {
			artifactName = artifactName[:18] + "..." + artifactName[len(artifactName)-19:]
		}
		shortSha := ""
		if len(artifact.Sha256) == 64 {
			shortSha = artifact.Sha256[:7] + "..." + artifact.Sha256[64-7:]
		}
		fmt.Printf(formatStringLine, "xxxx", artifactName, shortSha, since, len(artifact.CreationTimestamp))
	}

	return nil
}
