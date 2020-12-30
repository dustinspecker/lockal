#!/bin/bash
set -exo nounset
set -o pipefail

# clean up any previously installed executables
rm -rf ./bin

${LOCKAL_BIN} install

if [ ! -f ./bin/get_helm.sh ]; then
  echo "./bin/get_helm.sh was not created" >&2
  exit 1
fi

# create a modified get_helm.sh to cause the checksum to no longer match
echo "old file" > ./bin/get_helm.sh

${LOCKAL_BIN} install
actual_checksum="$(shasum -a 512 ./bin/get_helm.sh | awk '{ print $1 }')"
expected_checksum="$(cat lockal.star | awk '/checksum =/ { print $3 }' | sed 's/"//g' | sed 's/,//g')"

if [ "${actual_checksum}" != "${expected_checksum}" ]; then
  echo "checksum of ./bin/get_helm.sh did not match expected checksum" >&2
  exit 1
fi
