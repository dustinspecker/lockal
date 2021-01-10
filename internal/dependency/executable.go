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

type Executable struct {
	Name     string
	Location string
	Checksum string
}

func (exe Executable) Download(fs afero.Fs, logCtx *log.Entry, cacheDir string, getFile func(dest, src string) error) error {
	dest := fmt.Sprintf("bin/%s", exe.Name)

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
		removed, err := removeInvalidFile(fs, logCtx, dest, exe.Checksum)
		if err != nil {
			return err
		}

		if !removed {
			logCtx.Info(fmt.Sprintf("skipping download for %s as it already exists at %s", exe.Name, dest))
			return nil
		}

		logCtx.Info(fmt.Sprintf("removed old %s since it didn't match expected checksum", dest))
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

func downloadFile(fs afero.Fs, logCtx *log.Entry, location, dest, expectedChecksum string, getFile func(dest, src string) error) error {
	_, err := fs.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		logCtx.Info(fmt.Sprintf("downloading %s to %s", location, dest))

		if err := getFile(dest, location); err != nil {
			return err
		}

		removed, err := removeInvalidFile(fs, logCtx, dest, expectedChecksum)
		if err != nil {
			return err
		}

		if removed {
			errorMessage := fmt.Sprintf("downloaded %s did not match expected checksum", dest)
			logCtx.Error(errorMessage)

			return fmt.Errorf(errorMessage)
		}
	}

	return nil
}

func copyFile(fs afero.Fs, logCtx *log.Entry, src, dest string) error {
	logCtx.Info(fmt.Sprintf("copying from %s to %s", src, dest))

	if err := fs.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// copy from src to dest
	cacheContent, err := afero.ReadFile(fs, src)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, dest, cacheContent, 0755)
}

func markFileExecutable(fs afero.Fs, filepath string) error {
	return fs.Chmod(filepath, 0755)
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
