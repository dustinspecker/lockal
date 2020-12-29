package rules

import (
	"github.com/dustinspecker/lockal/internal/dependency"
	"go.starlark.net/starlark"
)

func Executable(addDep func(dep dependency.Dependency) error) func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var location string

		if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "name", &name, "location", &location); err != nil {
			return nil, err
		}

		addDep(dependency.Dependency{
			Name:     name,
			Location: location,
		})

		return starlark.None, nil
	}
}
