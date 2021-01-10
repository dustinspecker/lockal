package dependency

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/afero"

	"github.com/dustinspecker/lockal/internal/config"
)

type ExecutableFromArchive struct {
	Name               string
	Location           string
	ArchiveChecksum    string
	ExtractFilepath    string
	ExecutableChecksum string
}

func (efa ExecutableFromArchive) Download(cfg config.Config) error {
	dest := efa.Name

	// check if dest file exists
	// if dest file exists and checksum does match then do nothing
	// if dest file exists and checksum does not match, remove the old dest file
	// check cache for executable checksum, if not exist then extract from archive in cache
	//  -> check cache for archive checksum, if not exist then download new archive file to cache
	//	-> verify new archive file matches expected checksum, delete if no match
	// extract filepath from archive in cache to executable cache
	// copy executable file from cache to dest file
	// mark dest file as executable

	existingFileIsValid, err := validateExistingFile(cfg.Fs, cfg.LogCtx, dest, efa.ExecutableChecksum)
	if err != nil {
		return err
	}

	if existingFileIsValid {
		return nil
	}

	archiveCache := fmt.Sprintf("%s/lockal/sha512/%s/%s", cfg.CacheDir, efa.ArchiveChecksum[0:2], efa.ArchiveChecksum)
	if err = downloadFile(cfg.Fs, cfg.LogCtx, fmt.Sprintf("%s?archive=false", efa.Location), archiveCache, efa.ArchiveChecksum, cfg.GetFile); err != nil {
		return err
	}

	executableCache := fmt.Sprintf("%s/lockal/sha512/%s/%s", cfg.CacheDir, efa.ExecutableChecksum[0:2], efa.ExecutableChecksum)
	if err = extractFile(cfg.Fs, cfg.LogCtx, efa.Location, archiveCache, executableCache, efa.ExtractFilepath, efa.ExecutableChecksum, cfg.ExtractFileFromArchive); err != nil {
		return err
	}

	if err = copyFile(cfg.Fs, cfg.LogCtx, executableCache, dest); err != nil {
		return err
	}

	return markFileExecutable(cfg.Fs, dest)
}

func extractFile(fs afero.Fs, logCtx *log.Entry, archiveFileName, archiveCache, executableCache, extractFilepath, executableChecksum string, extractFileFromArchive func(archiveFileName, archivePath, extractFilepath, extractToDir string) error) error {
	// TODO: do nothing if executable file already exists in cache

	tempDir, err := afero.TempDir(fs, "", "")
	if err != nil {
		return err
	}

	logCtx.Info(fmt.Sprintf("extracting %s from %s to %s", extractFilepath, archiveCache, fmt.Sprintf("%s/%s", tempDir, extractFilepath)))

	if err := extractFileFromArchive(archiveFileName, archiveCache, extractFilepath, tempDir); err != nil {
		return err
	}

	if err := copyFile(fs, logCtx, fmt.Sprintf("%s/%s", tempDir, extractFilepath), executableCache); err != nil {
		return err
	}

	removed, err := removeInvalidFile(fs, logCtx, executableCache, executableChecksum)
	if err != nil {
		return err
	}

	if removed {
		errorMessage := fmt.Sprintf("extracted %s did not match expected checksum", extractFilepath)
		logCtx.Error(errorMessage)

		return fmt.Errorf(errorMessage)
	}

	return nil
}
