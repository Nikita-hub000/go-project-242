package main

import (
	"context"
	"fmt"
	"os"

	"code/internal/cliapp"
)

func main() {
	app := cliapp.New(os.Stdout)
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
