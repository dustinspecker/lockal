package parse

import (
	"github.com/spf13/afero"
	"go.starlark.net/starlark"

	"github.com/dustinspecker/lockal/internal/dependency"
	"github.com/dustinspecker/lockal/internal/rules"
)

func GetDependencies(fs afero.Fs) ([]dependency.Dependency, error) {
	deps := []dependency.Dependency{}

	fileData, err := fs.Open("lockal.star")
	if err != nil {
		return deps, err
	}

	thread := &starlark.Thread{
		Name: "lockal-main",
	}

	addDep := func(dep dependency.Dependency) error {
		deps = append(deps, dep)

		return nil
	}

	nativeFunctions := starlark.StringDict{
		"executable": starlark.NewBuiltin("executable", rules.Executable(addDep)),
	}

	_, err = starlark.ExecFile(thread, "lockal.star", fileData, nativeFunctions)

	return deps, err
}
