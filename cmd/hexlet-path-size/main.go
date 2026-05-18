package main

import (
	pathsize "code"
	"context"
	"fmt"
	"log"
	"os"

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
        Action: func(ctx context.Context, cmd *cli.Command) error {
            recursive := cmd.Bool("recursive")
            all := cmd.Bool("all")
            human := cmd.Bool("human")
			args := cmd.Args()
			if args.Len() == 0 {
				fmt.Println("path is required")
				return nil
			}
			path := args.First()
            info, err := pathsize.GetPathSize(path, human, all, recursive )            
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                return nil
            }
			if human {
				fmt.Printf("%s\t%s\n", info, path)
			} else {
				fmt.Printf("%s\t%s\n", info, path)
			}
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

