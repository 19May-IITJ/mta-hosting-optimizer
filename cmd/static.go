package main

import "github.com/urfave/cli/v2"

type appConfig struct {
	cli *cli.App
}

const (
	appname = "mta-hosting-optimizer"
	usage   = ""
	version = "0.0.1"
)
