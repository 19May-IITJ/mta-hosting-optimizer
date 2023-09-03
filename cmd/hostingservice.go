package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mta2/modules"
	"mta2/modules/utility"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			server := &http.Server{}

			utility.NATS_ADD = os.Getenv(utility.NATS_URI)
			port := os.Getenv(utility.HOSTINGSERVICE_PORT)

			if port != "" && utility.NATS_ADD != "" {
				err := modules.RegisterService(ctx, port, utility.HOSTINGSERVICE, server)
				if err != nil {
					return err
				}
			} else {
				//Logger Block
				log.Println("Unable to start Hosting Service")
				log.Println("Please provide ENV Variables or check ./mta help hosting service")
				return errors.New("unable to start Hosting Service please provide ENV Variables")
				//Logger Block
			}
			log.Println("Press ^C to shutdown server")
			<-sigChan
			// Handle OS signals (e.g., Ctrl+C).
			fmt.Println("\nReceived termination signal. Shutting down gracefully...")
			timeoutCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
			if err := server.Shutdown(timeoutCtx); err != nil {
				fmt.Printf("Server shutdown error: %v\n", err)
			}
			defer cancel()
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
