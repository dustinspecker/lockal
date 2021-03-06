package dependency

import (
	"fmt"
	"os"
	"testing"

	"github.com/apex/log"
	"github.com/spf13/afero"

	"github.com/dustinspecker/lockal/internal/config"
)

func TestDownload(t *testing.T) {
	// validate download without cache works
	fs := afero.NewMemMapFs()
	logHandler, logCtx := getLogCtx()

	getFile := func(dest, src string) error {
		if src != "some.sh/ghosthouse" {
			return fmt.Errorf("invalid src provided")
		}

		return afero.WriteFile(fs, dest, []byte("file a"), 0644)
	}

	exe := Executable{
		Name:     "bin/ghostdog",
		Location: "some.sh/ghosthouse",
		Checksum: "a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963",
	}

	cfg := config.Config{
		CacheDir: "/.cache",
		Fs:       fs,
		LogCtx:   logCtx,
		GetFile:  getFile,
	}

	err := exe.Download(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "downloading some.sh/ghosthouse to /.cache/lockal/sha512/a7/a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963") {
		t.Error("expected a log message saying download in progress")
	}

	stat, err := fs.Stat("bin/ghostdog")
	if err != nil {
		t.Fatalf("unexpected error stating bin/ghostdog: %v", err)
	}

	if stat.Mode() != 0755 {
		t.Errorf("expected file to be marked 0755, but was %v", stat.Mode())
	}

	// validate cache is used when possible
	if err = fs.Remove("bin/ghostdog"); err != nil {
		t.Fatalf("unexpected error while removing bin/ghostdog: %v", err)
	}

	getFileNoDownload := func(dest, src string) error {
		return fmt.Errorf("getFileNoDownload should not have been called - cache should have been used")
	}

	cfg.GetFile = getFileNoDownload

	err = exe.Download(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "copying from /.cache/lockal/sha512/a7/a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963 to bin/ghostdog") {
		t.Error("expected a log message saying download in progress")
	}

	stat, err = fs.Stat("bin/ghostdog")
	if err != nil {
		t.Fatalf("unexpected error stating bin/ghostdog: %v", err)
	}

	if stat.Mode() != 0755 {
		t.Errorf("expected file to be marked 0755, but was %v", stat.Mode())
	}

	cacheStat, err := fs.Stat("/.cache/lockal/sha512/a7/a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963")
	if err != nil {
		t.Fatalf("unexpected error stating cache of ghostdog: %v", err)
	}

	if cacheStat.Mode() == 0755 {
		t.Errorf("expected cache permissions to not be touched, but was marked executable")
	}
}

func TestDownloadSkipsGettingFileIfAlreadyExistsWithSameChecksum(t *testing.T) {
	fs := afero.NewMemMapFs()
	logHandler, logCtx := getLogCtx()

	if err := fs.Mkdir("bin", 0755); err != nil {
		t.Fatalf("unexpected error creating bin directory: %v", err)
	}

	if err := afero.WriteFile(fs, "bin/dustin", []byte("file dustin"), 0644); err != nil {
		t.Fatalf("unexpected error creating bin/dustin: %v", err)
	}

	getFile := func(dest, src string) error {
		t.Error("getFile should not have been called")

		return fmt.Errorf("should not be called")
	}

	exe := Executable{
		Name:     "bin/dustin",
		Location: "dustin.com/dustin",
		Checksum: "c301106040f367ce621cbafa373d73fe270a95aeb2a6076f15a6bf79c1634d39e67e62d3a660e410a865d1ec7e1c2a131270090083885656d1f941bdf8abefeb",
	}

	cfg := config.Config{
		CacheDir: "/home/dustin/.cache",
		Fs:       fs,
		LogCtx:   logCtx,
		GetFile:  getFile,
	}

	err := exe.Download(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "skipping download for bin/dustin as it already exists") {
		t.Error("expected a log message saying skipping download")
	}
}

func TestDownloadUpdatesExecutableIfChecksumMismatch(t *testing.T) {
	fs := afero.NewMemMapFs()
	logHandler, logCtx := getLogCtx()

	if err := fs.Mkdir("bin", 0755); err != nil {
		t.Fatalf("unexpected error creating bin directory: %v", err)
	}

	if err := afero.WriteFile(fs, "bin/dustin", []byte("old file dustin"), 0644); err != nil {
		t.Fatalf("unexpected error creating bin/dustin: %v", err)
	}

	getFile := func(dest, src string) error {
		return afero.WriteFile(fs, dest, []byte("file dustin"), 0644)
	}

	exe := Executable{
		Name:     "bin/dustin",
		Location: "dustin.com/dustin",
		Checksum: "c301106040f367ce621cbafa373d73fe270a95aeb2a6076f15a6bf79c1634d39e67e62d3a660e410a865d1ec7e1c2a131270090083885656d1f941bdf8abefeb",
	}

	cfg := config.Config{
		CacheDir: "/tmp/.cache",
		Fs:       fs,
		LogCtx:   logCtx,
		GetFile:  getFile,
	}

	err := exe.Download(cfg)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "removed old bin/dustin since it didn't match expected checksum") {
		t.Error("expected a log message saying updating old executable")
	}
}

func TestDownloadReturnsErrorIfChecksumDoesNotMatchAfterDownload(t *testing.T) {
	fs := afero.NewMemMapFs()
	logHandler, logCtx := getLogCtx()

	if err := fs.Mkdir("bin", 0755); err != nil {
		t.Fatalf("unexpected error creating bin directory: %v", err)
	}

	getFile := func(dest, src string) error {
		return afero.WriteFile(fs, dest, []byte("file a"), 0644)
	}

	exe := Executable{
		Name:     "ghostdog",
		Location: "some.sh/ghosthouse",
		Checksum: "hey",
	}

	cfg := config.Config{
		CacheDir: "/var/lib/lockal/.cache",
		Fs:       fs,
		LogCtx:   logCtx,
		GetFile:  getFile,
	}

	err := exe.Download(cfg)
	if err == nil {
		t.Fatal("expected an error when checksums do not match")
	}

	expectedErrorMessage := "downloaded /var/lib/lockal/.cache/lockal/sha512/he/hey did not match expected checksum"

	if err.Error() != expectedErrorMessage {
		t.Errorf("expected error message of \"%s\", but got \"%s\"", expectedErrorMessage, err.Error())
	}

	if !hasLogEntry(logHandler, log.ErrorLevel, log.Fields{"app": "lockal-test"}, expectedErrorMessage) {
		t.Error("expected a log message saying checksums did not match after download")
	}
	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "removing /var/lib/lockal/.cache/lockal/sha512/he/hey since it has a checksum of a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963, which does not match expected checksum of hey") {
		t.Error("expected a log message saying checksums did not match after download")
	}

	_, err = fs.Stat("bin/ghostdog")
	if err == nil {
		t.Fatal("expected an error stating bin/ghostdog not found")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected error to be IsNotExist, but got %v", err)
	}
}

func TestDownloadReturnsErrorWhenGetFileErrs(t *testing.T) {
	_, logCtx := getLogCtx()

	getFile := func(dest, src string) error {
		return fmt.Errorf("some error")
	}

	exe := Executable{
		Name:     "lockal",
		Location: "some.sh/lockal",
		Checksum: "samplechecksum",
	}

	cfg := config.Config{
		CacheDir: "/tmp/.cache",
		Fs:       afero.NewMemMapFs(),
		LogCtx:   logCtx,
		GetFile:  getFile,
	}

	err := exe.Download(cfg)
	if err == nil {
		t.Fatalf("expected error to be returned when getFile errs")
	}

	if err.Error() != "some error" {
		t.Errorf("expected error message of \"some error\" to be returned, but got \"%s\"", err.Error())
	}
}
