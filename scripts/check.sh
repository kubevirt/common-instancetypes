#!/bin/bash

# We need to add the current directory as a safe git directory before running any commands
if ! git config --global --get-all safe.directory | grep "$(pwd)" > /dev/null; then
    git config --global --add safe.directory "$(pwd)"
fi

if ! git diff --quiet common*.yaml README.md; then
    echo "Uncommited changes to generated files present, please commit these and try again."
    git --no-pager diff  common*.yaml README.md
    exit 1
fi