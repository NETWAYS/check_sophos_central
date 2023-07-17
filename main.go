package main

import (
	"github.com/NETWAYS/go-check"
)

const readme = `Check the status of alerts and endpoints over the API of the Sophos Central cloud service.

The plugin currently checks the state of all alerts and endpoints within a tenant, you need to supply API Token
(ID and secret) for a single tenant.`

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	plugin.Name = "check_sophos_central"
	plugin.Readme = readme
	plugin.Version = buildVersion()
	plugin.Timeout = 30

	// Parse arguments
	config := BuildConfigFlags(plugin.FlagSet)
	plugin.ParseArguments()
	config.SetFromEnv()

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	rc, output, err := config.Run()
	if err != nil {
		check.ExitError(err)
	}

	check.ExitRaw(rc, output)
}
