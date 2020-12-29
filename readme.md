# lockal

> a user/project local executable dependency manager

## Problems solved by lockal

I often setup new projects that require executables/binaries such as [shellcheck](https://www.shellcheck.net/) or
[kind](https://kind.sigs.k8s.io/). Since I have to support several projects it's easier for each project to have its own
install of these programs. This allows using multiple versions across the same machine, while enabling each project to use a specific version.
This process also prevents having to use root to install packages for a project. This also prevents one project needing a newer version that
would break builds for other projects, which is never wanted.

Unfortunately, this process of installing and managing executables results in me starting every new project with writing bash scripts and
make targets to install these.

Lockal automates this process.

## Install lockal

### download binary

### build from source

1. Have [go](https://golang.org/dl/) installed
1. Clone this repository
1. Navigate to the cloned repository
1. Run `make build`
1. `./bin/lockal` now exists in cloned repository
   - to make `lockal` available system-wide, run `sudo mv ./bin/lockal /usr/local/bin/lockal`

## Usage

Lockal uses the [Starlark language](https://github.com/bazelbuild/starlark), specifically the [go implementation](https://github.com/google/starlark-go).
Starlark gives us the power of a scripting language similar to Python for configuration, while preventing non-determinstic results such as file access.

### lockal.star

Let's start with an example where we want to download a bash script.

An example `lockal.star` to install a bash script from https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3
and name the file `get_helm.sh`.

```starlark
executable(
  name = "get_helm.sh",
  location = "https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3",
)
```

In the directory where `lockal.star` exists (typically the project root), run
`lockal install`. Lockal will analyze the `lockal.star` file and begin downloading
executables.

We'll see lockal log something similar to:

```
2020/12/29 08:34:31 downloading get_helm.sh from https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3 to bin/get_helm.sh
```

Lockal will have created the bash script at `./bin/get_helm.sh` (relative to where `lockal.star` exists.

If we execute `lockal install` again, we'll see the following logs:

```
2020/12/29 08:37:32 skipping download for get_helm.sh as it already exists at bin/get_helm.sh
```

Lockal only downloads what doesn't exist.

## Commands

### `lockal install`

`lockal install` ensures each executable defined in `lockal.star` is installed.
