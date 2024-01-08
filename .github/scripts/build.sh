#!/bin/bash

set -ex

if [ -z "${COMMON_INSTANCETYPES_VERSION}" ] ; then
    echo "${BASH_SOURCE[0]} expects the following env variables to be provided: COMMON_INSTANCETYPES_VERSION."
    exit 1
fi

# Build the bundles
COMMON_INSTANCETYPES_VERSION=${COMMON_INSTANCETYPES_VERSION} make

cd _build || exit 1
for file in common-*-bundle.yaml; do
    file_versioned=${file/\.yaml/-${COMMON_INSTANCETYPES_VERSION}\.yaml}
    mv "${file}" "${file_versioned}"
    sed -i "s/${file}/${file_versioned}/g" CHECKSUMS.sha256
done
