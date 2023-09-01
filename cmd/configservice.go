package main

import (
	"context"
	"errors"
	"log"
	"mta2/modules"
	"mta2/modules/utility"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func configservice(a *appConfig) *cli.Command {
	return &cli.Command{
		//EventRouter
		Name:      utility.CONFIGSERVICE,
		Usage:     utility.CONFIGSERVICEUSAGE,
		ArgsUsage: "",
		Before:    beforeConfigService,
		Action: func(c *cli.Context) error {
			utility.NATS_ADD = os.Getenv(utility.NATS_URI)
			port := os.Getenv(utility.CONFIGSERVICE_PORT)
			if port != "" && utility.NATS_ADD != "" {
				modules.RegisterService(port, utility.CONFIGSERVICE)
			} else {
				//Logger Block
				log.Println("Unable to start Config Service")
				log.Println("Please provide ENV Variables or check ./mta help configservice")
				return errors.New("unable to start Config Service please provide ENV Variables")
				//Logger Block
			}
			log.Println("Press ^C to shutdown server")
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			<-ctx.Done()
			defer stop()
			return nil

		},
		After: func(ctx *cli.Context) error {
			return nil
		},
	}
}
func beforeConfigService(c *cli.Context) error {

	return nil
}
func configservice_help() {

}
