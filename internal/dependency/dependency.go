package dependency

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"

	"github.com/apex/log"

	"github.com/spf13/afero"
)

type Dependency struct {
	Name     string
	Location string
	Checksum string
}

type GetFile func(dest, src string) error

func (dep Dependency) Download(fs afero.Fs, logCtx *log.Entry, getFile GetFile) error {
	dest := fmt.Sprintf("bin/%s", dep.Name)

	_, err := fs.Stat(dest)

	// check if file exists
	// if file exists and checksum does not match, remove the old file and download the new file
	// if file exists and checksum does match then do nothing
	// if file does not exist then download new file
	//	-> verify new file matches expected checksum, delete if no match

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		removed, err := removeInvalidFile(fs, logCtx, dest, dep.Checksum)
		if err != nil {
			return err
		}

		if !removed {
			logCtx.Info(fmt.Sprintf("skipping download for %s as it already exists at %s", dep.Name, dest))
			return nil
		}

		logCtx.Info(fmt.Sprintf("removed old %s since it didn't match expected checksum", dest))
	}

	logCtx.Info(fmt.Sprintf("downloading %s from %s to %s", dep.Name, dep.Location, dest))

	if err := getFile(dest, dep.Location); err != nil {
		return err
	}

	removed, err := removeInvalidFile(fs, logCtx, dest, dep.Checksum)
	if err != nil {
		return err
	}

	if removed {
		errorMessage := fmt.Sprintf("downloaded %s did not match expected checksum", dest)
		logCtx.Error(errorMessage)

		return fmt.Errorf(errorMessage)
	}

	return fs.Chmod(dest, 0755)
}

func removeInvalidFile(fs afero.Fs, logCtx *log.Entry, targetPath, expectedChecksum string) (bool, error) {
	actualChecksum, err := getChecksum(fs, targetPath)
	if err != nil {
		return false, err
	}

	// checksum matches, so don't remove file
	if actualChecksum == expectedChecksum {
		return false, nil
	}

	logCtx.Info(fmt.Sprintf("removing %s since it has a checksum of %s, which does not match expected checksum of %s", targetPath, actualChecksum, expectedChecksum))

	if err = fs.Remove(targetPath); err != nil {
		return false, err
	}

	return true, nil
}

func getChecksum(fs afero.Fs, filepath string) (string, error) {
	fileContent, err := fs.Open(filepath)
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	hash := sha512.New()

	if _, err = io.Copy(hash, fileContent); err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
