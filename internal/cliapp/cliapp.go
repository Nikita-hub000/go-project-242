package cliapp

import (
	"context"
	"io"

	"code"
	"code/internal/sizefmt"

	"github.com/urfave/cli/v3"
)

const pathArgumentMessage = "exactly one path argument is required"

func New(stdout io.Writer) *cli.Command {
	return &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory",
		ArgsUsage: "<path>",
		Writer:    stdout,
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
				return cli.Exit(pathArgumentMessage, 1)
			}

			path := args.First()
			size, err := code.GetPathSize(path, cmd.Bool("recursive"), cmd.Bool("human"), cmd.Bool("all"))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if _, err := io.WriteString(cmd.Root().Writer, sizefmt.FormatLine(size, path)); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
