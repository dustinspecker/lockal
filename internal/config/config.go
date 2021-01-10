package config

import (
	"github.com/apex/log"
	"github.com/spf13/afero"
)

type Config struct {
	CacheDir               string
	Fs                     afero.Fs
	LogCtx                 *log.Entry
	GetFile                func(dest, src string) error
	ExtractFileFromArchive func(archiveFileName, archivePath, extractFilepath, extractToDir string) error
}
