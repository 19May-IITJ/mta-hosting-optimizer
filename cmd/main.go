package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	rootctx := context.Background()
	_, cancel := context.WithCancel(rootctx)
	defer cancel()
	app := buildCLI()
	_ = app.cli.Run(os.Args)

}

func buildCLI() *appConfig {
	return builds(cli.NewApp())
}
