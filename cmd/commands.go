package main

import (
	"fmt"
	"mta2/modules/utility"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func commandsForExecution_HostingService(a *appConfig) *cli.Command {
	return hostingservice(a)
}
func commandsForExecution_ConfigService(a *appConfig) *cli.Command {
	return configservice(a)
}
func help() *cli.Command {
	return &cli.Command{

		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "helper function",
		Action: func(ctx *cli.Context) error {
			a := ctx.Args().Get(0)
			switch a {
			case utility.HOSTINGSERVICE:
				hostingservice_help()
			case utility.CONFIGSERVICE:
				configservice_help()
			default:
				fmt.Println()
				fmt.Printf("Usage : ./mta COMMANDS ")
				fmt.Println(usage)
				headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
				columnFmt := color.New(color.FgYellow).SprintfFunc()
				fmt.Println()
				fmt.Println("Commands:")
				tblcommands := table.New("CommandName", "Usage")
				tblcommands.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
				tblcommands.AddRow(utility.HOSTINGSERVICE, utility.HOSTINGSERVICEUSAAGE)
				tblcommands.AddRow(utility.CONFIGSERVICE, utility.CONFIGSERVICEUSAGE)
				tblcommands.AddRow("help", "provide console help for application or command")
				tblcommands.Print()
			}
			return cli.Exit("Application is stopped", 0)
		},
	}
}
