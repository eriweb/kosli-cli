package main

import (
	"fmt"
	"testing"

	"github.com/kosli-dev/cli/internal/requests"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type AllowArtifactCommandTestSuite struct {
	suite.Suite
	defaultKosliArguments string
	envName               string
	artifactName          string
	// artifactPath          string
	// fingerprint string
}

func (suite *AllowArtifactCommandTestSuite) SetupTest() {
	suite.envName = "allow-artifact-env"
	suite.artifactName = "arti"
	// suite.artifactPath = "testdata/folder1/hello.txt"
	global = &GlobalOpts{
		ApiToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6ImNkNzg4OTg5In0.e8i_lA_QrEhFncb05Xw6E_tkCHU9QfcY4OLTVUCHffY",
		Owner:    "docs-cmd-test-user",
		Host:     "http://localhost:8001",
	}
	kosliClient = requests.NewKosliClient(1, false, logger)
	suite.defaultKosliArguments = fmt.Sprintf(" --host %s --owner %s --api-token %s", global.Host, global.Owner, global.ApiToken)

	// CreateFlow(suite.flowName, suite.T())
	// fingerprintOptions := &fingerprintOptions{
	// 	artifactType: "file",
	// }
	// var err error
	// suite.fingerprint, err = GetSha256Digest(suite.artifactPath, fingerprintOptions, logger)
	// require.NoError(suite.T(), err)
	// CreateArtifact(suite.flowName, suite.fingerprint, suite.artifactName, suite.T())
	CreateEnv(global.Owner, suite.envName, "server", suite.T())
}

func (suite *AllowArtifactCommandTestSuite) TestAllowArtifactCmd() {
	tests := []cmdTestCase{
		{
			name:   "allowing an artifact works with --fingerprint",
			cmd:    fmt.Sprintf(`allow artifact %s --fingerprint 8e568bd886069f1290def0caabc1e97ce0e7b80c105e611258b57d76fcef234c  --environment %s --reason because %s`, suite.artifactName, suite.envName, suite.defaultKosliArguments),
			golden: "artifact 8e568bd886069f1290def0caabc1e97ce0e7b80c105e611258b57d76fcef234c was allow listed in environment: allow-artifact-env\n",
		},
		{
			name:   "allowing an artifact works with --artifact-type",
			cmd:    fmt.Sprintf(`allow artifact testdata/file1  --artifact-type file  --environment %s --reason because %s`, suite.envName, suite.defaultKosliArguments),
			golden: "artifact 7509e5bda0c762d2bac7f90d758b5b2263fa01ccbc542ab5e3df163be08e6ca9 was allow listed in environment: allow-artifact-env\n",
		},
		{
			wantError: true,
			name:      "allowing an artifact fails if artifact name argument is missing",
			cmd:       fmt.Sprintf(`allow artifact --fingerprint 8e568bd886069f1290def0caabc1e97ce0e7b80c105e611258b57d76fcef234c  --environment %s --reason because %s`, suite.envName, suite.defaultKosliArguments),
			golden:    "Error: accepts 1 arg(s), received 0\n",
		},
		{
			wantError: true,
			name:      "allowing an artifact fails if --reason is missing",
			cmd:       fmt.Sprintf(`allow artifact %s --fingerprint 8e568bd886069f1290def0caabc1e97ce0e7b80c105e611258b57d76fcef234c  --environment %s %s`, suite.artifactName, suite.envName, suite.defaultKosliArguments),
			golden:    "Error: required flag(s) \"reason\" not set\n",
		},
		{
			wantError: true,
			name:      "allowing an artifact fails if --environment is missing",
			cmd:       fmt.Sprintf(`allow artifact %s --fingerprint 8e568bd886069f1290def0caabc1e97ce0e7b80c105e611258b57d76fcef234c  --reason because %s`, suite.artifactName, suite.defaultKosliArguments),
			golden:    "Error: required flag(s) \"environment\" not set\n",
		},
		{
			wantError: true,
			name:      "allowing an artifact fails if --fingerprint and --artifact-type are missing",
			cmd:       fmt.Sprintf(`allow artifact %s --environment %s --reason because %s`, suite.artifactName, suite.envName, suite.defaultKosliArguments),
			golden:    "Error: either --artifact-type or --fingerprint must be specified\nUsage: kosli allow artifact {IMAGE-NAME | FILE-PATH | DIR-PATH} [flags]\n",
		},
	}

	runTestCmd(suite.T(), tests)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAllowArtifactCommandTestSuite(t *testing.T) {
	suite.Run(t, new(AllowArtifactCommandTestSuite))
}
