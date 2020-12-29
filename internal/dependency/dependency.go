package dependency

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
)

type Dependency struct {
	Name     string
	Location string
}

type GetFile func(dest, src string) error

func (dep Dependency) Download(fs afero.Fs, getFile GetFile) error {
	dest := fmt.Sprintf("bin/%s", dep.Name)

	_, err := fs.Stat(dest)

	if err == nil {
		// file exists
		log.Printf("skipping download for %s as it already exists at %s", dep.Name, dest)

		return nil
	}

	if !os.IsNotExist(err) {
		// permission error? etc. trying to read file
		return err
	}

	// file does not exist, so download it
	log.Printf("downloading %s from %s to %s", dep.Name, dep.Location, dest)

	if err := getFile(dest, dep.Location); err != nil {
		return err
	}

	return fs.Chmod(dest, 0755)
}
