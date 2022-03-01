---
title: "merkely pipeline artifact report evidence bitbucket-pullrequest"
---

## merkely pipeline artifact report evidence bitbucket-pullrequest

Report a Bitbucket pull request evidence for an artifact in a Merkely pipeline.

### Synopsis


   Check if a pull request exists for an artifact and report the pull-request evidence to the artifact in Merkely. 
   The artifact SHA256 fingerprint is calculated or alternatively it can be provided directly. 
   The following flags are defaulted as follows in the CI list below:

   
	| Bitbucket 
	|---------------------------------------------------------------------------
	| build-url : https://bitbucket.org/${BITBUCKET_WORKSPACE}/${BITBUCKET_REPO_SLUG}/addon/pipelines/home#!/results/${BITBUCKET_BUILD_NUMBER}
	|---------------------------------------------------------------------------
	| Github 
	|---------------------------------------------------------------------------
	| build-url : ${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}
	|---------------------------------------------------------------------------
	| Teamcity 
	|---------------------------------------------------------------------------
	|---------------------------------------------------------------------------

```shell
merkely pipeline artifact report evidence bitbucket-pullrequest [ARTIFACT-NAME-OR-PATH] [flags]
```

### Options

```
  -t, --artifact-type string         The type of the artifact to calculate its SHA256 fingerprint. One of: [docker, file, dir]
      --assert                       Fail if no pull requests found for the given commit.
      --bitbucket-password string    Bitbucket password.
      --bitbucket-username string    Bitbucket user name.
      --bitbucket-workspace string   Bitbucket workspace.
  -b, --build-url string             The url of CI pipeline that generated the evidence. (default "https://github.com/merkely-development/cli/actions/runs/1915357107")
      --commit string                Git commit for which to find pull request evidence. (default "db7e69d466d9958ec7ce9574f3c5b042ef903f63")
  -d, --description string           [optional] The evidence description.
  -e, --evidence-type string         The type of evidence being reported.
  -h, --help                         help for bitbucket-pullrequest
  -p, --pipeline string              The Merkely pipeline name.
      --registry-password string     The docker registry password or access token.
      --registry-provider string     The docker registry provider or url.
      --registry-username string     The docker registry username.
      --repository string            Git repository. (default "merkely-development/cli")
  -s, --sha256 string                The SHA256 fingerprint for the artifact. Only required if you don't specify --artifact-type.
```

### Options inherited from parent commands

```
  -a, --api-token string      The merkely API token.
  -c, --config-file string    [optional] The merkely config file path. (default "merkely")
  -D, --dry-run               Whether to run in dry-run mode. When enabled, data is not sent to Merkely and the CLI exits with 0 exit code regardless of errors.
  -H, --host string           The merkely endpoint. (default "https://app.merkely.com")
  -r, --max-api-retries int   How many times should API calls be retried when the API host is not reachable. (default 3)
  -o, --owner string          The merkely user or organization.
  -v, --verbose               Print verbose logs to stdout.
```
