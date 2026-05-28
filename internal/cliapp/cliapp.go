package cliapp

import (
	"context"
	"fmt"
	"io"

	pathsize "code"

	"github.com/urfave/cli/v3"
)

const pathArgumentError = "exactly one path argument is required"

func New(stdout io.Writer) *cli.Command {
	return &cli.Command{
		Name:           "hexlet-path-size",
		Usage:          "print size of a file or directory",
		ArgsUsage:      "<path>",
		Writer:         stdout,
		ExitErrHandler: returnExitError,
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
			args := cmd.Args()
			if args.Len() != 1 {
				return cli.Exit(pathArgumentError, 1)
			}

			path := args.First()
			size, err := pathsize.GetPathSize(path, cmd.Bool("recursive"), cmd.Bool("human"), cmd.Bool("all"))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			_, err = fmt.Fprintf(cmd.Root().Writer, "%s\t%s\n", size, path)
			return err
		},
	}
}

func returnExitError(_ context.Context, _ *cli.Command, _ error) {}
