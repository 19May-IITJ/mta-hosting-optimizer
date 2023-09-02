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

func configservice(a *appConfig) *cli.Command {
	return &cli.Command{
		//EventRouter
		Name:      utility.CONFIGSERVICE,
		Usage:     utility.CONFIGSERVICEUSAGE,
		ArgsUsage: "",
		Before:    beforeConfigService,
		Action: func(c *cli.Context) error {

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			server := &http.Server{}

			utility.NATS_ADD = os.Getenv(utility.NATS_URI)
			port := os.Getenv(utility.CONFIGSERVICE_PORT)
			if port != "" && utility.NATS_ADD != "" {
				modules.RegisterService(ctx, port, utility.CONFIGSERVICE, server)
			} else {
				//Logger Block
				log.Println("Unable to start Config Service")
				log.Println("Please provide ENV Variables or check ./mta help configservice")
				return errors.New("unable to start Config Service please provide ENV Variables")
				//Logger Block
			}
			log.Println("Press ^C to shutdown server")
			// Wait for signals or context cancellation.

			<-sigChan
			// Handle OS signals (e.g., Ctrl+C).
			fmt.Println("\nReceived termination signal. Shutting down gracefully...")
			timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
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
func beforeConfigService(c *cli.Context) error {

	return nil
}
func configservice_help() {

}
