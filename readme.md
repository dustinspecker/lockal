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

Once this is all scripted, then the scenario of how to handle a developer having a stale executable appears. This
is cumbersome and typically results in attempting to parse the output of some version command to determine if the executable is the expected
version. Then we have to handle removing stale executables and downloading the correct one.

Also, sometimes downloading is slow, so it would be great if Lockal could access a cache of previous executables it has downloaded.

Then a coworker says, "Hey, your scripts are cool, but they only support Linux. I use a Mac."

One morning you wake up to see someone on the other side of the planet created a GitHub issue stating, "I want to use this on my Raspberry Pi.
It's Linux, but arm64 architecture."

Some projects release plain binaries, but some projects release archive files containing executables. So, we need a way to handle extracting a file from
an archive.

Lockal helps automate this process.

## Install lockal

### download binary

1. Download the [latest binary](https://github.com/dustinspecker/lockal/releases/latest) for your operating system and architecture.
2. Mark the binary as exectuable by using `chmod`
   - for example, `chmod +x lockal-linux-amd64`
3. To make `lockal` available system-wide, move it to `/usr/local/bin`
   - for example, `sudo mv lockal-linux-amd64 /usr/local/bin/lockal`

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

### Download a single executable

Let's start with an example where we want to download a bash script. We'll first create a `lockal.star` file to define desired dependencies.

An example `lockal.star` to install a bash script from https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3
and name the file `get_helm.sh` looks like:

```starlark
executable(
  name = "bin/get_helm.sh",
  location = "https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3",
  checksum = "6faf31a30425399b7d75ad2d00cfcca12725b0386387b5569f382d6f7aecf123996c11f5d892c74236face3801d511dd9f1ec52e744ad3adfb397269f4c0c2bc",
)
```

`name` is the name of the file that should be created, so a file named `get_helm.sh` in a directory named `bin` will be created. `location` is where to download
the file from. Lockal requires `checksum` to determine if it should update a stale executable. Lockal also uses the `checksum`
to validate that the expected artifact was downloaded.

> Note: checksum *must* be a sha512

In the directory where `lockal.star` exists (typically the project root), run
`lockal install`. Lockal will analyze the `lockal.star` file and begin downloading
executables.

Lockal will log that it has downloaded `get_helm.sh` to its cache and created the bash script at `bin/get_helm.sh` (relative to where `lockal.star` exists).

If we execute `lockal install` again, Lockal will log that it's skipping `get_helm.sh`.

> Note: Lockal only downloads executables that don't exist or when the expected checksum does not match what exists.

### Download multiple executables

Rarely, do we only need one executable for a project, so let's add another to our `lockal.star` file:

```starlark
executable(
  name = "bin/get_helm.sh",
  location = "https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3",
  checksum = "6faf31a30425399b7d75ad2d00cfcca12725b0386387b5569f382d6f7aecf123996c11f5d892c74236face3801d511dd9f1ec52e744ad3adfb397269f4c0c2bc",
)

executable(
  name = "bin/kind",
  location = "https://github.com/kubernetes-sigs/kind/releases/download/v0.9.0/kind-linux-amd64",
  checksum = "e7152acf5fd7a4a56af825bda64b1b8343a1f91588f9b3ddd5420ae5c5a95577d87431f2e417a7e03dd23914e1da9bed855ec19d0c4602729b311baccb30bd7f",
)
```

`lockal install` will skip `get_helm.sh` since it already exists, but will retrieve `kind`.

### Download specific executables based on operating system and architecture

Our current `lockal.star` only supports Linux/amd64, but we can also support Mac/amd64 and Linux/arm64 with a few changes:

```starlark
executable(
  name = "bin/get_helm.sh",
  location = "https://raw.githubusercontent.com/helm/helm/23dd3af5e19a02d4f4baa5b2f242645a1a3af629/scripts/get-helm-3",
  checksum = "6faf31a30425399b7d75ad2d00cfcca12725b0386387b5569f382d6f7aecf123996c11f5d892c74236face3801d511dd9f1ec52e744ad3adfb397269f4c0c2bc",
)

def get_kind_checksum(os, arch):
  if os == "linux":
    if arch == "amd64":
      return "e7152acf5fd7a4a56af825bda64b1b8343a1f91588f9b3ddd5420ae5c5a95577d87431f2e417a7e03dd23914e1da9bed855ec19d0c4602729b311baccb30bd7f"
    if arch == "arm64":
      return "0bcb81fe7e3aa4515df0c3c7607b3cd6f3cf2e87b029f18b4c4628e15225062d543cd1abfc8ac56477f159177f16fab4e416d598dc1beb57ad8ed46e9e6b180d"

  if os == "darwin":
    return "1b716be0c6371f831718bb9f7e502533eb993d3648f26cf97ab47c2fa18f55c7442330bba62ba822ec11edb84071ab616696470cbdbc41895f2ae9319a7e3a99"

  fail("unsupported operating_system/architecture: %s/%s" % (os, arch))

executable(
  name = "bin/kind",
  location = "https://github.com/kubernetes-sigs/kind/releases/download/v0.9.0/kind-%(os)s-%(arch)s" % dict(os = LOCKAL_OS, arch = LOCKAL_ARCH),
  checksum = get_kind_checksum(LOCKAL_OS, LOCKAL_ARCH),
)
```

Now `lockal install` will retrieve the `kind` executable for Linux or Mac (darwin) as desired.

### Download and extract an executable from an archive

It's common for projects to release artifacts in an archive such as a `tar.gz` file. Lockal can also handle this.

An example `lockal.star`:

```starlark
executable_from_archive(
  name = "bin/helm",
  location = "https://get.helm.sh/helm-v3.4.2-linux-amd64.tar.gz",
  archive_checksum = "f827744743df68c11f619f64f0f7c915c1afd15673ee287c5b8d68cf3c246deae97ac86aadd761e22432d7b5e927fc65288ce3dca80a495af6b2aefa71bce22a",
  extract_filepath = "linux-amd64/helm",
  executable_checksum = "d89093f1c463355b7280017c357a7d86825548a96d6b6772ae07fcc76a25474d02d3ba8f125514c49ff83383410863cd8b56702c5f9dcfa1f3f0d23ac1587fa1",
)
```

Notice we now use `executable_from_archive` instead of `executable`. The `name` behaves the same
and `location` is now the location of the archive to download.

`archive_checksum` is used to cache and validate the downloaded archive matches what is expected.
`extract_filepath` is a filepath within the archive to retrieve. `executable_checksum` is then
used to cache and validate the extracted executable is what's expected.

For this example, Lockal will
1. download the archive
2. validate the downloaded archive against the `archive_checksum`
3. extract `linux-amd64/helm` from the archive
4. validate the extacted file against the `executable_checksum`
5. place the extracted file in `bin/helm`

## Commands

### `lockal install`

`lockal install` ensures each executable defined in `lockal.star` is installed.

### `lockal version`

`lockal version` prints the version of Lockal being used
