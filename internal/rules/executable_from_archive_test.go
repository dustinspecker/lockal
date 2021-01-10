package rules

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/lockal/internal/dependency"
)

func TestExecutableFromArchive(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("some_efa_name"),
		starlark.String("some_efa_location"),
		starlark.String("some_efa_archive_checksum"),
		starlark.String("some_efa_extract_filepath"),
		starlark.String("some_efa_executable_checksum"),
	}
	kwargs := []starlark.Tuple{}

	addDepCalled := false

	addDep := func(dep dependency.Dependency) error {
		addDepCalled = true

		efa := dep.(dependency.ExecutableFromArchive)

		if efa.Name != "some_efa_name" {
			t.Errorf("expected dep.Name to be some_efa_name, but was %s", efa.Name)
		}

		if efa.Location != "some_efa_location" {
			t.Errorf("expected efa.Location to be some_efa_location, but was %s", efa.Location)
		}

		if efa.ArchiveChecksum != "some_efa_archive_checksum" {
			t.Errorf("expected efa.ArchiveChecksum to be some_efa_archive_checksum, but was %s", efa.ArchiveChecksum)
		}

		if efa.ExtractFilepath != "some_efa_extract_filepath" {
			t.Errorf("expected efa.ExtractFilepath to be some_efa_extract_filepath, but was %s", efa.ExtractFilepath)
		}

		if efa.ExecutableChecksum != "some_efa_executable_checksum" {
			t.Errorf("expected efa.ExecutableChecksum to be some_efa_executable_checksum, but was %s", efa.ExecutableChecksum)
		}

		return nil
	}

	value, err := ExecutableFromArchive(addDep)(thread, builtin, args, kwargs)
	if err != nil {
		t.Fatalf("unexpected error invoking ExecutableFromArchive: %v", err)
	}

	if value != starlark.None {
		t.Errorf("expected value to be None, but got: %v", value)
	}

	if !addDepCalled {
		t.Error("expected addDep to be called")
	}
}

func TestExecutableFromArchiveReturnsErrorWhenInvalidArgs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{}
	kwargs := []starlark.Tuple{}

	addDep := func(dep dependency.Dependency) error {
		return nil
	}

	_, err := ExecutableFromArchive(addDep)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("ExecutableFromArchive should have returned an error")
	}
}
