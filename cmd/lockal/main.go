package main

import (
	"os"

	"github.com/apex/log"
	cliHandler "github.com/apex/log/handlers/cli"
	gogetter "github.com/hashicorp/go-getter"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/lockal/internal/parse"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetHandler(cliHandler.New(os.Stderr))

	logCtx := log.WithFields(log.Fields{
		"app": "lockal",
	})

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
						if err = dep.Download(afero.NewOsFs(), logCtx, getFile); err != nil {
							return err
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logCtx.Fatal(err.Error())
	}
}
