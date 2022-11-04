#!/bin/bash

# We need to add the current directory as a safe git directory before running any commands
git config --global --add safe.directory "$(pwd)"

git diff --quiet common*.yaml README.md
uncommited=$?
if [[ $uncommited -eq 1 ]]; then
    echo "Uncommited changes to generated files present, please commit these and try again."
    git --no-pager diff  common*.yaml README.md
    exit 1
fi