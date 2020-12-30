package dependency

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/spf13/afero"
)

func TestDownload(t *testing.T) {
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
	}

	err := dep.Download(fs, logCtx, getFile)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "downloading ghostdog from some.sh/ghosthouse to bin/ghostdog") {
		t.Error("expected a log message saying download in progress")
	}

	stat, err := fs.Stat("bin/ghostdog")
	if err != nil {
		t.Fatalf("unexpected error stating bin/ghostdog: %v", err)
	}

	if stat.Mode() != 0755 {
		t.Errorf("expected file to be marked 0755, but was %v", stat.Mode())
	}
}

func TestDownloadSkipsGettingFileIfAlreadyExists(t *testing.T) {
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
	}

	err := dep.Download(fs, logCtx, getFile)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if !hasLogEntry(logHandler, log.InfoLevel, log.Fields{"app": "lockal-test"}, "skipping download for dustin as it already exists at bin/dustin") {
		t.Error("expected a log message saying skipping download")
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
	}

	err := dep.Download(afero.NewMemMapFs(), logCtx, getFile)
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
