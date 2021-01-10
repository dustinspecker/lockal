package dependency

import (
	"fmt"

	"github.com/dustinspecker/lockal/internal/config"
)

type Executable struct {
	Name     string
	Location string
	Checksum string
}

func (exe Executable) Download(cfg config.Config) error {
	dest := exe.Name

	// check if dest file exists
	// if dest file exists and checksum does match then do nothing
	// if dest file exists and checksum does not match, remove the old dest file
	// check cache for checksum, if not exist then download new file to cache
	//	-> verify new file matches expected checksum, delete if no match
	// copy file from cache to dest file
	// mark dest file as executable

	existingFileIsValid, err := validateExistingFile(cfg.Fs, cfg.LogCtx, dest, exe.Checksum)
	if err != nil {
		return err
	}

	if existingFileIsValid {
		return nil
	}

	cache := fmt.Sprintf("%s/lockal/sha512/%s/%s", cfg.CacheDir, exe.Checksum[0:2], exe.Checksum)
	if err = downloadFile(cfg.Fs, cfg.LogCtx, exe.Location, cache, exe.Checksum, cfg.GetFile); err != nil {
		return err
	}

	if err = copyFile(cfg.Fs, cfg.LogCtx, cache, dest); err != nil {
		return err
	}

	return markFileExecutable(cfg.Fs, dest)
}
