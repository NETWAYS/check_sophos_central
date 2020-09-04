package main

import (
	"github.com/NETWAYS/go-check"
)

const readme = `Check the status of alerts and endpoints over the API of the Sophos Central cloud service.

The plugin currently checks the state of all alerts and endpoints within a tenant, you need to supply API Token
(ID and secret) for a single tenant.

https://github.com/NETWAYS/check_sophos_central

Copyright (c) 2020 NETWAYS GmbH <info@netways.de>
Copyright (c) 2020 Markus Frosch <markus.frosch@netways.de

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see https://www.gnu.org/licenses/`

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

	check.Exit(rc, output)
}
