#!/bin/bash
set -ex

for test_script in test/integration/*.sh ; do
  "${test_script}"
done
