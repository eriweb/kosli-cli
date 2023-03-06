package main

import (
	"io"

	ghUtils "github.com/kosli-dev/cli/internal/github"
	"github.com/spf13/cobra"
)

type assertPullRequestGithubOptions struct {
	githubConfig *ghUtils.GithubConfig
	commit       string
}

const assertPRGithubShortDesc = `Assert if a Github pull request for a git commit exists.`

const assertPRGithubLongDesc = assertPRGithubShortDesc + `
The command exits with non-zero exit code 
if no pull requests were found for the commit.`

const assertPRGithubExample = `
kosli assert github-pullrequest  \
	--github-token yourGithubToken \
	--github-org yourGithubOrg \
	--commit yourArtifactGitCommit \
	--commit yourGitCommit \
	--repository yourGithubGitRepository
`

func newAssertPullRequestGithubCmd(out io.Writer) *cobra.Command {
	o := new(assertPullRequestGithubOptions)
	o.githubConfig = new(ghUtils.GithubConfig)
	cmd := &cobra.Command{
		Use:     "github-pullrequest",
		Aliases: []string{"gh-pr", "github-pr"},
		Short:   assertPRGithubShortDesc,
		Long:    assertPRGithubLongDesc,
		Example: assertPRGithubExample,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(args)
		},
	}

	ci := WhichCI()
	addGithubFlags(cmd, o.githubConfig, ci)
	cmd.Flags().StringVar(&o.commit, "commit", DefaultValue(ci, "git-commit"), commitPREvidenceFlag)
	addDryRunFlag(cmd)

	err := RequireFlags(cmd, []string{"github-token", "github-org", "commit", "repository"})
	if err != nil {
		logger.Error("failed to configure required flags: %v", err)
	}

	return cmd
}

func (o *assertPullRequestGithubOptions) run(args []string) error {
	// repository name must be extracted if a user is using default value from ${GITHUB_REPOSITORY}
	// because the value is in the format of "owner/repository"
	o.githubConfig.Repository = extractRepoName(o.githubConfig.Repository)
	pullRequestsEvidence, err := getPullRequestsEvidence(o.githubConfig, o.commit, true)
	if err != nil {
		return err
	}
	logger.Info("found [%d] pull request(s) in Github for commit: %s", len(pullRequestsEvidence), o.commit)
	return nil
}
