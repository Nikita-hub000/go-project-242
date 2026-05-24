package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"code/internal/pathsize"
	"code/internal/sizefmt"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			recursive := cmd.Bool("recursive")
			all := cmd.Bool("all")
			human := cmd.Bool("human")
			args := cmd.Args()

			if args.Len() == 0 {
				return cli.Exit("path is required", 1)
			}

			path := args.First()

			size, err := pathsize.Calculate(path, all, recursive)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			info := sizefmt.BytesRaw(size)
			if human {
				info = sizefmt.BytesIEC(size)
			}

			fmt.Printf("%s\t%s\n", info, path)
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		var exitErr cli.ExitCoder
		if errors.As(err, &exitErr) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
