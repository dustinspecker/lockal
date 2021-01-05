package rules

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/lockal/internal/dependency"
)

func TestExecutable(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("some_name"),
		starlark.String("some_location"),
		starlark.String("some_checksum"),
	}
	kwargs := []starlark.Tuple{}

	addDepCalled := false

	addDep := func(exe dependency.Executable) error {
		addDepCalled = true

		if exe.Name != "some_name" {
			t.Errorf("expected dep.Name to be some_name, but was %s", exe.Name)
		}

		if exe.Location != "some_location" {
			t.Errorf("expected exe.Location to be some_location, but was %s", exe.Location)
		}

		if exe.Checksum != "some_checksum" {
			t.Errorf("expected exe.Checksum to be some_checksum, but was %s", exe.Checksum)
		}

		return nil
	}

	value, err := Executable(addDep)(thread, builtin, args, kwargs)
	if err != nil {
		t.Fatalf("unexpected error invoking Executable: %v", err)
	}

	if value != starlark.None {
		t.Errorf("expected value to be None, but got: %v", value)
	}

	if !addDepCalled {
		t.Error("expected addDep to be called")
	}
}

func TestExecutableReturnsErrorWhenInvalidArgs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{}
	kwargs := []starlark.Tuple{}

	addDep := func(exe dependency.Executable) error {
		return nil
	}

	_, err := Executable(addDep)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("Files should have returned an error")
	}
}
