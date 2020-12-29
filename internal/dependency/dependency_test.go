package dependency

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
)

func TestDownload(t *testing.T) {
	fs := afero.NewMemMapFs()

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

	err := dep.Download(fs, getFile)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
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

	err := dep.Download(fs, getFile)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestDownloadReturnsErrorWhenGetFileErrs(t *testing.T) {
	getFile := func(dest, src string) error {
		return fmt.Errorf("some error")
	}

	dep := Dependency{
		Name:     "lockal",
		Location: "some.sh/lockal",
	}

	err := dep.Download(afero.NewMemMapFs(), getFile)
	if err == nil {
		t.Fatalf("expected error to be returned when getFile errs")
	}

	if err.Error() != "some error" {
		t.Errorf("expected error message of \"some error\" to be returned, but got \"%s\"", err.Error())
	}
}
