package dependency

import (
	"fmt"
	"os"

	"github.com/apex/log"

	"github.com/spf13/afero"
)

type Dependency struct {
	Name     string
	Location string
}

type GetFile func(dest, src string) error

func (dep Dependency) Download(fs afero.Fs, logCtx *log.Entry, getFile GetFile) error {
	dest := fmt.Sprintf("bin/%s", dep.Name)

	_, err := fs.Stat(dest)

	if err == nil {
		// file exists
		logCtx.Info(fmt.Sprintf("skipping download for %s as it already exists at %s", dep.Name, dest))

		return nil
	}

	if !os.IsNotExist(err) {
		// permission error? etc. trying to read file
		return err
	}

	// file does not exist, so download it
	logCtx.Info(fmt.Sprintf("downloading %s from %s to %s", dep.Name, dep.Location, dest))

	if err := getFile(dest, dep.Location); err != nil {
		return err
	}

	return fs.Chmod(dest, 0755)
}
