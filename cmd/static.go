package main

import "github.com/urfave/cli/v2"

type appConfig struct {
	cli *cli.App
}

const (
	appname = "mta-hosting-optimizer"
	usage   = "uncovers the inefficient servers hosting only few active MTAs depending on threshold value by user"
	version = "0.0.1"
)
