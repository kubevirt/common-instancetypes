#!/bin/bash

set -e

# Basic instancetype.kubevirt.io version check
INSTANCETYPE_API_VERSION="v1beta1"
for dir in instancetypes preferences; do
    if grep -riP "instancetype.kubevirt.io/v\d+" "${dir}" | grep -v ${INSTANCETYPE_API_VERSION}; then
        echo "Unexpected instancetype.kubevirt.io version detected in tree! Expected instancetype.kubevirt.io/${INSTANCETYPE_API_VERSION}, please fix or update scripts/check.sh"
        exit 1
    fi
done

# We need to add the current directory as a safe git directory before running any commands
if ! git config --global --get-all safe.directory | grep "$(pwd)" > /dev/null; then
    git config --global --add safe.directory "$(pwd)"
fi

if ! git diff --quiet README.md; then
    echo "Uncommited changes to generated files present, please commit these and try again."
    git --no-pager diff README.md
    exit 1
fi
