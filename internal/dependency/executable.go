package dependency

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/afero"
)

type Executable struct {
	Name     string
	Location string
	Checksum string
}

func (exe Executable) Download(fs afero.Fs, logCtx *log.Entry, cacheDir string, getFile func(dest, src string) error) error {
	dest := exe.Name

	// check if dest file exists
	// if dest file exists and checksum does match then do nothing
	// if dest file exists and checksum does not match, remove the old dest file
	// check cache for checksum, if not exist then download new file to cache
	//	-> verify new file matches expected checksum, delete if no match
	// copy file from cache to dest file
	// mark dest file as executable

	existingFileIsValid, err := validateExistingFile(fs, logCtx, dest, exe.Checksum)
	if err != nil {
		return err
	}

	if existingFileIsValid {
		return nil
	}

	cache := fmt.Sprintf("%s/lockal/sha512/%s/%s", cacheDir, exe.Checksum[0:2], exe.Checksum)
	if err = downloadFile(fs, logCtx, exe.Location, cache, exe.Checksum, getFile); err != nil {
		return err
	}

	if err = copyFile(fs, logCtx, cache, dest); err != nil {
		return err
	}

	return markFileExecutable(fs, dest)
}
