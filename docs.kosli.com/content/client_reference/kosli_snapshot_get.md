---
title: "kosli snapshot get"
---

## kosli snapshot get

Get a specific environment snapshot.

### Synopsis

Get a specific environment snapshot.

```shell
kosli snapshot get ENVIRONMENT-NAME-OR-EXPRESSION [flags]
```

### Flags
| Flag | Description |
| :--- | :--- |
|    -h, --help  |  help for get  |
|    -j, --json  |  [optional] Print output as json.  |


### Options inherited from parent commands
| Flag | Description |
| :--- | :--- |
|    -a, --api-token string  |  The Kosli API token.  |
|    -c, --config-file string  |  [optional] The Kosli config file path. (default "kosli")  |
|    -D, --dry-run  |  [optional] Whether to run in dry-run mode. When enabled, data is not sent to Kosli and the CLI exits with 0 exit code regardless of errors.  |
|    -H, --host string  |  [defaulted] The Kosli endpoint. (default "https://app.kosli.com")  |
|    -r, --max-api-retries int  |  [defaulted] How many times should API calls be retried when the API host is not reachable. (default 3)  |
|        --owner string  |  The Kosli user or organization.  |
|    -v, --verbose  |  [optional] Print verbose logs to stdout.  |


### Examples

```shell

# get the latest snapshot of an environment:
kosli snapshot get yourEnvironmentName
	--api-token yourAPIToken \
	--owner yourOrgName 

# get the SECOND latest snapshot of an environment:
kosli snapshot get yourEnvironmentName~1
	--api-token yourAPIToken \
	--owner yourOrgName 

# get the snapshot number 23 of an environment:
kosli snapshot get yourEnvironmentName#23
	--api-token yourAPIToken \
	--owner yourOrgName 

```
