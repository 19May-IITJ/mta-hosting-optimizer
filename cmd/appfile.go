package main

import (
	"github.com/urfave/cli/v2"
)

func builds(app *cli.App) *appConfig {

	a := appConfig{}
	a.cli = app

	a.cli.Name = appname
	a.cli.Usage = usage
	a.cli.Version = version

	a.cli.Commands = []*cli.Command{
		commandsForExecution_HostingService(&a),
		commandsForExecution_ConfigService(&a),
		help(),
	}
	return &a
}
