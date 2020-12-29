package parse

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGetDependency(t *testing.T) {
	fs := afero.NewMemMapFs()

	fileContents := `
executable(
	name = "cat",
	location = "farm/feline",
)

executable(
	name = "cloud",
	location = "sky/cloud",
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
	if deps[1].Name != "cloud" {
		t.Errorf("expected second dep to have name cloud, but got %s", deps[1].Name)
	}
	if deps[1].Location != "sky/cloud" {
		t.Errorf("expected second dep to have location sky/cloud, but got %s", deps[1].Location)
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
