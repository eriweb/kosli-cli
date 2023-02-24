package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kosli-dev/cli/internal/output"
	"github.com/kosli-dev/cli/internal/requests"
	"github.com/spf13/cobra"
)

const logEnvironmentShortDesc = `List environment events.`

const logEnvironmentLongDesc = logEnvironmentShortDesc + `
The results are paginated and ordered from latests to oldest. 
By default, the page limit is 15 events per page.

You can optionally specify an INTERVAL between two snapshot expressions with <expression>..<expression>.
Expressions can be:
	~N   N'th behind the latest snapshot
	N    snapshot number N
	NOW  the latest snapshot
Either expression can be omitted to default to NOW.
`

const logEnvironmentExample = `
# list the last 15 events for an environment:
kosli log environment yourEnvironmentName \
	--api-token yourAPIToken \
	--owner yourOrgName

# list the last 30 events for an environment:
kosli log environment yourEnvironmentName \
	--page-limit 30 \
	--api-token yourAPIToken \
	--owner yourOrgName

# list the last 30 events for an environment (in JSON):
kosli log environment yourEnvironmentName \
	--page-limit 30 \
	--api-token yourAPIToken \
	--owner yourOrgName \
	--output json
`

type logEnvironmentOptions struct {
	listOptions
	reverse bool
}

func newLogEnvironmentCmd(out io.Writer) *cobra.Command {
	o := new(logEnvironmentOptions)
	cmd := &cobra.Command{
		Use:     "environment ENV_NAME [INTERVAL]",
		Aliases: []string{"env"},
		Short:   logEnvironmentShortDesc,
		Long:    logEnvironmentLongDesc,
		Example: logEnvironmentExample,
		Args:    cobra.MatchAll(cobra.MaximumNArgs(2), cobra.MinimumNArgs(1)),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Owner", "ApiToken"})
			if err != nil {
				return ErrorBeforePrintingUsage(cmd, err.Error())
			}

			return o.validate(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out, args)
		},
	}

	addListFlags(cmd, &o.listOptions)
	cmd.Flags().BoolVar(&o.reverse, "reverse", false, reverseFlag)

	return cmd
}

func (o *logEnvironmentOptions) run(out io.Writer, args []string) error {
	envName := args[0]
	interval := ""
	if len(args) > 1 {
		interval = args[1]
	}

	return o.getEnvironmentEvents(out, envName, interval)

}

// events

func (o *logEnvironmentOptions) getEnvironmentEvents(out io.Writer, envName, interval string) error {
	url := fmt.Sprintf("%s/api/v1/environments/%s/%s/events/?page=%d&per_page=%d&interval=%s&reverse=%t",
		global.Host, global.Owner, envName, o.pageNumber, o.pageLimit, url.QueryEscape(interval), o.reverse)

	reqParams := &requests.RequestParams{
		Method:   http.MethodGet,
		URL:      url,
		Password: global.ApiToken,
	}
	response, err := kosliClient.Do(reqParams)
	if err != nil {
		return err
	}
	return output.FormattedPrint(response.Body, o.output, out, o.pageNumber,
		map[string]output.FormatOutputFunc{
			"table": printEnvironmentEventsLogAsTable,
			"json":  output.PrintJson,
		})
}
