package dependency

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apex/log"

	"github.com/spf13/afero"
)

type Dependency struct {
	Name     string
	Location string
	Checksum string
}

type GetFile func(dest, src string) error

func (dep Dependency) Download(fs afero.Fs, logCtx *log.Entry, cacheDir string, getFile GetFile) error {
	dest := fmt.Sprintf("bin/%s", dep.Name)

	_, err := fs.Stat(dest)

	// check if dest file exists
	// if dest file exists and checksum does match then do nothing
	// if dest file exists and checksum does not match, remove the old dest file
	// check cache for checksum, if not exist then download new file to cache
	//	-> verify new file matches expected checksum, delete if no match
	// copy file from cache to dest file
	// mark dest file as executable

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

	cache := fmt.Sprintf("%s/lockal/sha512/%s/%s", cacheDir, dep.Checksum[0:2], dep.Checksum)
	_, err = fs.Stat(cache)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		logCtx.Info(fmt.Sprintf("downloading %s from %s to %s", dep.Name, dep.Location, cache))

		if err := getFile(cache, dep.Location); err != nil {
			return err
		}

		removed, err := removeInvalidFile(fs, logCtx, cache, dep.Checksum)
		if err != nil {
			return err
		}

		if removed {
			errorMessage := fmt.Sprintf("downloaded %s did not match expected checksum", cache)
			logCtx.Error(errorMessage)

			return fmt.Errorf(errorMessage)
		}

		if err = fs.Chmod(cache, 0755); err != nil {
			return err
		}
	}

	logCtx.Info(fmt.Sprintf("copying from %s to %s", cache, dest))

	if err = fs.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// copy from cache to dest
	cacheContent, err := afero.ReadFile(fs, cache)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, dest, cacheContent, 0755)
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
