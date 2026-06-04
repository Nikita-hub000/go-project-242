package main

import (
	"context"
	"os"

	"code/internal/cliapp"
)

func main() {
	app := cliapp.New(os.Stdout)
	app.ErrWriter = os.Stderr

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
