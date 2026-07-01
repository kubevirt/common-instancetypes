#!/bin/bash

set -e

# We need to add the current directory as a safe git directory before running any commands
if ! git config --global --get-all safe.directory | grep "$(pwd)" > /dev/null; then
    git config --global --add safe.directory "$(pwd)"
fi

if ! git diff --quiet README.md tests; then
    echo "Uncommited changes to generated files present, please commit these and try again."
    git --no-pager diff README.md tests
    exit 1
fi
