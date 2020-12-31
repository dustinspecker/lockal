#!/bin/bash
set -exo nounset
set -o pipefail

# clean up any previously installed executables
rm -rf ./bin

${LOCKAL_BIN} install

if [ ! -f ./bin/kind ]; then
  echo "./bin/kind was not created" >&2
  exit 1
fi
