package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"code/internal/cliapp"

	"github.com/urfave/cli/v3"
)

func main() {
	app := cliapp.New(os.Stdout)
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode(err))
	}
}

func exitCode(err error) int {
	var exitCoder cli.ExitCoder
	if errors.As(err, &exitCoder) {
		return exitCoder.ExitCode()
	}

	return 1
}
