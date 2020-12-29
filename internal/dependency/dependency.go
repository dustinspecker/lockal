package dependency

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
)

type Dependency struct {
	Name     string
	Location string
}

type GetFile func(dest, src string) error

func (dep Dependency) Download(fs afero.Fs, getFile GetFile) error {
	dest := fmt.Sprintf("bin/%s", dep.Name)

	log.Printf("downloading %s from %s to %s", dep.Name, dep.Location, dest)

	if err := getFile(dest, dep.Location); err != nil {
		return err
	}

	if err := fs.Chmod(dest, 0755); err != nil {
		return err
	}

	return nil
}
