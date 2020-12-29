#!/bin/bash
set -exo nounset

# clean up any previously installed executables
rm -rf ./bin

${LOCKAL_BIN} install

if [ ! -f ./bin/get_helm.sh ]; then
  echo "./bin/get_helm.sh was not created" >&2
  exit 1
fi
