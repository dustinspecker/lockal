package parse

import (
	"runtime"

	"github.com/spf13/afero"
	"go.starlark.net/starlark"

	"github.com/dustinspecker/lockal/internal/dependency"
	"github.com/dustinspecker/lockal/internal/rules"
)

func GetDependencies(fs afero.Fs) ([]dependency.Executable, error) {
	deps := []dependency.Executable{}

	fileData, err := fs.Open("lockal.star")
	if err != nil {
		return deps, err
	}

	thread := &starlark.Thread{
		Name: "lockal-main",
	}

	addDep := func(exe dependency.Executable) error {
		deps = append(deps, exe)

		return nil
	}

	nativeFunctions := starlark.StringDict{
		"LOCKAL_ARCH": starlark.String(runtime.GOARCH),
		"LOCKAL_OS":   starlark.String(runtime.GOOS),
		"executable":  starlark.NewBuiltin("executable", rules.Executable(addDep)),
	}

	_, err = starlark.ExecFile(thread, "lockal.star", fileData, nativeFunctions)

	return deps, err
}
