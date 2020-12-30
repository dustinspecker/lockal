package dependency

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/spf13/afero"
)

func TestDownload(t *testing.T) {
	// validate download without cache works
	fs := afero.NewMemMapFs()
	logHandler, logCtx := getLogCtx()

	if err := fs.Mkdir("bin", 0755); err != nil {
		t.Fatalf("unexpected error creating bin directory: %v", err)
	}

	getFile := func(dest, src string) error {
		if src != "some.sh/ghosthouse" {
			return fmt.Errorf("invalid src provided")
		}

		return afero.WriteFile(fs, dest, []byte("file a"), 0644)
	}

	dep := Dependency{
		Name:     "ghostdog",
		Location: "some.sh/ghosthouse",
		Checksum: "a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963",
	}

	err := dep.Download(fs, logCtx, "/.cache", getFile)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "downloading ghostdog from some.sh/ghosthouse to /.cache/lockal/sha512/a7/a705aaf587ddc9ed135d4c318c339f3a0d6eb3a2e11936942afbfcd65254da6a1600b7b8e27f59464219fdc704f3b96c9953d80c05632411f475eea6f4548963") {
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

	err = dep.Download(fs, logCtx, "/.cache", getFileNoDownload)
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

	dep := Dependency{
		Name:     "dustin",
		Location: "dustin.com/dustin",
		Checksum: "c301106040f367ce621cbafa373d73fe270a95aeb2a6076f15a6bf79c1634d39e67e62d3a660e410a865d1ec7e1c2a131270090083885656d1f941bdf8abefeb",
	}

	err := dep.Download(fs, logCtx, "/home/dustin/.cache", getFile)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "skipping download for dustin as it already exists at bin/dustin") {
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

	dep := Dependency{
		Name:     "dustin",
		Location: "dustin.com/dustin",
		Checksum: "c301106040f367ce621cbafa373d73fe270a95aeb2a6076f15a6bf79c1634d39e67e62d3a660e410a865d1ec7e1c2a131270090083885656d1f941bdf8abefeb",
	}

	err := dep.Download(fs, logCtx, "/tmp/.cache", getFile)
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

	dep := Dependency{
		Name:     "ghostdog",
		Location: "some.sh/ghosthouse",
		Checksum: "hey",
	}

	err := dep.Download(fs, logCtx, "/var/lib/lockal/.cache", getFile)
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

	dep := Dependency{
		Name:     "lockal",
		Location: "some.sh/lockal",
		Checksum: "samplechecksum",
	}

	err := dep.Download(afero.NewMemMapFs(), logCtx, "/tmp/.cache", getFile)
	if err == nil {
		t.Fatalf("expected error to be returned when getFile errs")
	}

	if err.Error() != "some error" {
		t.Errorf("expected error message of \"some error\" to be returned, but got \"%s\"", err.Error())
	}
}

func getLogCtx() (*memory.Handler, *log.Entry) {
	log.SetLevel(log.DebugLevel)
	logHandler := memory.New()
	log.SetHandler(logHandler)

	logCtx := log.WithFields(log.Fields{
		"app": "lockal-test",
	})

	return logHandler, logCtx
}

func hasLogEntry(handler *memory.Handler, logLevel log.Level, fields log.Fields, message string) bool {
	for _, entry := range handler.Entries {
		if entry.Level == logLevel && entry.Message == message && reflect.DeepEqual(entry.Fields, fields) {
			return true
		}
	}

	return false
}
