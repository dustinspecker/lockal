#!/bin/bash
set -ex

export LOCKAL_BIN=$PWD/bin/lockal

# verify each example works
for test_file in $(find ./examples/ -name test.sh -type file); do
  echo $test_file
  set +e

  pushd $(dirname $test_file)
  ./test.sh
  rc=$?
  popd

  set -e

  if [ $rc -eq 0 ]; then
    echo "$test_file passed"
  else
    echo "$test_file failed"
    exit 1
  fi
done
