#!/bin/bash
set -ex

expected_version=$(git describe --always --dirty --long --tags)

actual_version="$($PWD/bin/lockal version)"

if [ "${actual_version}" != "${expected_version}" ]; then
  echo "expected version to be ${expected_version}, but got ${actual_version}" >&2
  exit 1
fi
