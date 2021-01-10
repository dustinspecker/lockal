kind_location = "https://github.com/kubernetes-sigs/kind/releases/download/v0.9.0/kind-%(os)s-%(arch)s" % dict(
  os = LOCKAL_OS,
  arch = LOCKAL_ARCH,
)

def get_kind_checksum(operating_system, cpu_architecture):
  if operating_system == "linux":
    if cpu_architecture == "amd64":
      return "e7152acf5fd7a4a56af825bda64b1b8343a1f91588f9b3ddd5420ae5c5a95577d87431f2e417a7e03dd23914e1da9bed855ec19d0c4602729b311baccb30bd7f"
    if cpu_architecture == "arm64":
      return "0bcb81fe7e3aa4515df0c3c7607b3cd6f3cf2e87b029f18b4c4628e15225062d543cd1abfc8ac56477f159177f16fab4e416d598dc1beb57ad8ed46e9e6b180d"

  if operating_system == "darwin":
    return "1b716be0c6371f831718bb9f7e502533eb993d3648f26cf97ab47c2fa18f55c7442330bba62ba822ec11edb84071ab616696470cbdbc41895f2ae9319a7e3a99"

  fail("unsupported operating_system/cpu_architecture: %s/%s" % (operating_system, cpu_architecture))

executable(
  name = "bin/kind",
  location = kind_location,
  checksum = get_kind_checksum(LOCKAL_OS, LOCKAL_ARCH),
)
