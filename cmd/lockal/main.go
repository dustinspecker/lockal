package main

import (
	"log"
	"os"

	gogetter "github.com/hashicorp/go-getter"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/lockal/internal/parse"
)

func main() {
	app := &cli.App{
		Name:  "lockal",
		Usage: "manage binary dependencies",
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "install dependencies from lockal.star",
				Action: func(c *cli.Context) error {
					deps, err := parse.GetDependencies(afero.NewOsFs())
					if err != nil {
						return err
					}

					getFile := func(dest, src string) error {
						return gogetter.GetFile(dest, src)
					}
					for _, dep := range deps {
						if err = dep.Download(afero.NewOsFs(), getFile); err != nil {
							return err
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
