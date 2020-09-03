package main

import (
	"github.com/NETWAYS/go-check"
)

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	plugin.Name = "check_sophos_central"
	plugin.Readme = `Check the status of alerts and endpoints over the API of the Sophos Central cloud service`
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

	check.Exit(rc, output)
}
