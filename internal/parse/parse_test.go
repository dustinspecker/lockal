package parse

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/spf13/afero"
)

func TestGetDependency(t *testing.T) {
	fs := afero.NewMemMapFs()

	fileContents := `
executable(
	name = "cat",
	location = "farm/feline",
	checksum = "some_sum",
)

executable(
	name = "cloud",
	location = "sky/cloud-%(os)s-%(arch)s" % dict(os = LOCKAL_OS, arch = LOCKAL_ARCH),
	checksum = "another_sum",
)
`

	if err := afero.WriteFile(fs, "lockal.star", []byte(fileContents), 0644); err != nil {
		t.Fatalf("unexpected error while creating lockal.star: %v", err)
	}

	deps, err := GetDependencies(fs)
	if err != nil {
		t.Fatalf("unexpected error when invoking GetDependencies: %v", err)
	}

	if len(deps) != 2 {
		t.Fatalf("expected 2 deps to be returned, but got %d", len(deps))
	}

	if deps[0].Name != "cat" {
		t.Errorf("expected first dep to have name cat, but got %s", deps[0].Name)
	}
	if deps[0].Location != "farm/feline" {
		t.Errorf("expected first dep to have location farm/feline, but got %s", deps[0].Location)
	}
	if deps[0].Checksum != "some_sum" {
		t.Errorf("expected first dep to have checksum some_sum, but got %s", deps[0].Checksum)
	}
	if deps[1].Name != "cloud" {
		t.Errorf("expected second dep to have name cloud, but got %s", deps[1].Name)
	}
	expectedDepLocation := fmt.Sprintf("sky/cloud-%s-%s", runtime.GOOS, runtime.GOARCH)
	if deps[1].Location != expectedDepLocation {
		t.Errorf("expected second dep to have location %s, but got %s", expectedDepLocation, deps[1].Location)
	}
	if deps[1].Checksum != "another_sum" {
		t.Errorf("expected second dep to have checksum another_sum, but got %s", deps[1].Checksum)
	}
}

func TestGetDependencyReturnsErrorIfLockalFileNotFound(t *testing.T) {
	_, err := GetDependencies(afero.NewMemMapFs())
	if err == nil {
		t.Fatalf("expected error when lockal.star file not found, but got no error")
	}
}

func TestGetDependencyReturnsErrorIfLockalFileIsInvalid(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, "lockal.star", []byte("not_a_rule()"), 0644); err != nil {
		t.Fatalf("unexpected error while creating lockal.star: %v", err)
	}

	_, err := GetDependencies(fs)
	if err == nil {
		t.Fatalf("expected error when lockal.star file is not valid")
	}
}
