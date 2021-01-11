package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	cliHandler "github.com/apex/log/handlers/cli"
	gogetter "github.com/hashicorp/go-getter"
	"github.com/mholt/archiver"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/lockal/internal/config"
	"github.com/dustinspecker/lockal/internal/parse"
)

var (
	VERSION = "dev"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetHandler(cliHandler.New(os.Stderr))

	logCtx := log.WithFields(log.Fields{
		"app": "lockal",
	})

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		logCtx.WithError(err).Fatal("getting home directory")
	}

	app := &cli.App{
		Name:  "lockal",
		Usage: "manage binary dependencies",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "log-level",
				Usage:  "level of logs to write (debug, info, warn, error, fatal)",
				Value:  "info",
				Hidden: false,
			},
		},
		Before: func(c *cli.Context) error {
			logLevel, err := log.ParseLevel(c.String("log-level"))
			if err != nil {
				return err
			}

			log.SetLevel(logLevel)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "install dependencies from lockal.star",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cache-directory",
						Usage:   "where to save cached downloads",
						Value:   filepath.Join(userHomeDir, ".cache"),
						EnvVars: []string{"XDG_CACHE_DIR"},
					},
				},
				Action: func(c *cli.Context) error {
					deps, err := parse.GetDependencies(afero.NewOsFs())
					if err != nil {
						return err
					}

					getFile := func(dest, src string) error {
						return gogetter.GetFile(dest, src)
					}

					extractFileFromArchive := func(archiveFileName, archivePath, extractFilepath, extractToDir string) error {
						extractorType, err := archiver.ByExtension(archiveFileName)
						if err != nil {
							return err
						}

						extractor, ok := extractorType.(archiver.Extractor)
						if !ok {
							return fmt.Errorf("invalid extractor")
						}

						return extractor.Extract(archivePath, extractFilepath, extractToDir)
					}

					cfg := config.Config{
						CacheDir:               c.String("cache-directory"),
						Fs:                     afero.NewOsFs(),
						LogCtx:                 logCtx,
						GetFile:                getFile,
						ExtractFileFromArchive: extractFileFromArchive,
					}

					for _, dep := range deps {
						if err = dep.Download(cfg); err != nil {
							return err
						}
					}

					return nil
				},
			},
			{
				Name:  "version",
				Usage: "print version of lockal",
				Action: func(c *cli.Context) error {
					fmt.Println(VERSION)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logCtx.Fatal(err.Error())
	}
}
