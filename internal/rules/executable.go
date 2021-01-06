package rules

import (
	"github.com/dustinspecker/lockal/internal/dependency"
	"go.starlark.net/starlark"
)

func Executable(addDep func(dep dependency.Dependency) error) func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var location string
		var checksum string

		if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "name", &name, "location", &location, "checksum", &checksum); err != nil {
			return nil, err
		}

		addDep(dependency.Executable{
			Name:     name,
			Location: location,
			Checksum: checksum,
		})

		return starlark.None, nil
	}
}
