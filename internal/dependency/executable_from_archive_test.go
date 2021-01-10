package dependency

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"

	"github.com/dustinspecker/lockal/internal/config"
)

func TestExecutableFromArchiveDownload(t *testing.T) {
	fs := afero.NewMemMapFs()
	_, logCtx := getLogCtx()

	efa := ExecutableFromArchive{
		Name:               "exe",
		Location:           "http://archive.tgz",
		ArchiveChecksum:    "21b9c6c34401c466769ec75e894d47f3d5eb656358ae836dc6d87b7747af69377f8266913427dfcd0027e68873ae8962f8afd943a29ccfacacabd27113a981be",
		ExtractFilepath:    "artifacts/executable",
		ExecutableChecksum: "bc07ffe5b4dbd2c52c87bce5298893c63e38a0d0333e2e01bbcfeddfdd40602724400d2998cb2a75e216aaffc913306a908d6057729a76102086b19556dc8be2",
	}

	getFile := func(dest, src string) error {
		if src == "http://archive.tgz?archive=false" {
			return afero.WriteFile(fs, dest, []byte("an archive"), 0644)
		}

		return nil
	}

	extractFileFromArchive := func(archiveFileName, archivePath, extractFilepath, extractToDir string) error {
		if archiveFileName != "http://archive.tgz" {
			t.Errorf("expected archive file to be http://archive.tgz, but got %s", archiveFileName)
		}

		expectedArchivePath := "/.cache/lockal/sha512/21/21b9c6c34401c466769ec75e894d47f3d5eb656358ae836dc6d87b7747af69377f8266913427dfcd0027e68873ae8962f8afd943a29ccfacacabd27113a981be"
		if archivePath != expectedArchivePath {
			t.Errorf("expected archivePath to be %s, but got %s", expectedArchivePath, archivePath)
		}

		dest := fmt.Sprintf("%s/%s", extractToDir, extractFilepath)

		return afero.WriteFile(fs, dest, []byte("an executable"), 0644)
	}

	cfg := config.Config{
		CacheDir:               "/.cache",
		Fs:                     fs,
		LogCtx:                 logCtx,
		GetFile:                getFile,
		ExtractFileFromArchive: extractFileFromArchive,
	}

	if err := efa.Download(cfg); err != nil {
		t.Fatalf("unexpected error when invoking Download: %v", err)
	}

	getFileShouldNotBeCalled := func(dest, src string) error {
		return fmt.Errorf("getFile should not be called when archive exists in cache")
	}

	extractFileFromArchiveShouldNotBeCalled := func(archiveFileName, archivePath, extractFilepath, extractToDir string) error {
		return fmt.Errorf("extractFileFromArchive should not be called when extracted file exists in cache")
	}

	cfg.GetFile = getFileShouldNotBeCalled
	cfg.ExtractFileFromArchive = extractFileFromArchiveShouldNotBeCalled

	if err := efa.Download(cfg); err != nil {
		t.Fatalf("unexpected error when invoking Download after cache populated: %v", err)
	}
}
