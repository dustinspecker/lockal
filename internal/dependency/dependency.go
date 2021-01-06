package dependency

import (
	"github.com/apex/log"
	"github.com/spf13/afero"
)

type Dependency interface {
	Download(fs afero.Fs, logCtx *log.Entry, cacheDir string, getFile func(string, string) error) error
}
