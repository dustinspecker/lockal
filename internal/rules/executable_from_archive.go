package rules

import (
	"github.com/dustinspecker/lockal/internal/dependency"
	"go.starlark.net/starlark"
)

func ExecutableFromArchive(addDep func(dep dependency.Dependency) error) func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var location string
		var archiveChecksum string
		var extractFilepath string
		var executableChecksum string

		if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "name", &name, "location", &location, "archive_checksum", &archiveChecksum, "extract_filepath", &extractFilepath, "executable_checksum", &executableChecksum); err != nil {
			return nil, err
		}

		addDep(dependency.ExecutableFromArchive{
			Name:               name,
			Location:           location,
			ArchiveChecksum:    archiveChecksum,
			ExtractFilepath:    extractFilepath,
			ExecutableChecksum: executableChecksum,
		})

		return starlark.None, nil
	}
}
