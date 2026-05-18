package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)
func GetPathSize(path string, all bool, recursive bool) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	answer := int64(0)

	for _, file := range files {
		if !all && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(path, file.Name())

		if file.IsDir() && !recursive {
			continue
		}

		size, err := GetPathSize(fullPath, all, recursive)
		if err != nil {
			return 0, err
		}

		answer += size
	}

	return answer, nil
}
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
            info, err := GetPathSize(path, all, recursive)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                return nil
            }
			if human {
				fmt.Printf("%s\t%s\n", ByteCountIEC(info), path)
			} else {
				fmt.Printf("%dB\t%s\n", info, path)
			}
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func ByteCountIEC(info int64) string {
	switch {
	case info < 1024:
		return fmt.Sprintf("%dB", info)
	case info < 1024*1024:
		return fmt.Sprintf("%.1fKB", float64(info)/1024)
	case info < 1024*1024*1024:
		return fmt.Sprintf("%.1fMB", float64(info)/(1024*1024))
	default:
		return fmt.Sprintf("%.1fGB", float64(info)/(1024*1024*1024))
	}
}
