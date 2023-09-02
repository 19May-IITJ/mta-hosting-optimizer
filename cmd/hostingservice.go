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

func hostingservice(a *appConfig) *cli.Command {
	return &cli.Command{
		//EventRouter
		Name:      utility.HOSTINGSERVICE,
		Usage:     utility.HOSTINGSERVICE,
		ArgsUsage: "",
		Before:    beforeHostingService,
		Action: func(c *cli.Context) error {
			{
				os.Setenv(utility.HOSTINGSERVICE_PORT, "8040")
				os.Setenv(utility.NATS_URI, "nats://localhost:4222")
			}
			utility.NATS_ADD = os.Getenv(utility.NATS_URI)
			port := os.Getenv(utility.HOSTINGSERVICE_PORT)

			if port != "" && utility.NATS_ADD != "" {
				modules.RegisterService(port, utility.HOSTINGSERVICE)
			} else {
				//Logger Block
				log.Println("Unable to start Hosting Service")
				log.Println("Please provide ENV Variables or check ./mta help hosting service")
				return errors.New("unable to start Hosting Service please provide ENV Variables")
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

func beforeHostingService(c *cli.Context) error {

	return nil
}
func hostingservice_help() {

}
