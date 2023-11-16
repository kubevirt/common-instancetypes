#!/bin/bash

set -ex

if [ -z "${COMMON_INSTANCETYPES_VERSION}" ] ; then
    echo "${BASH_SOURCE[0]} expects the following env variables to be provided: COMMON_INSTANCETYPES_VERSION."
    exit 1
fi

# Build the bundles
COMMON_INSTANCETYPES_VERSION=${COMMON_INSTANCETYPES_VERSION} make

cd _build || exit 1
for f in common-*-bundle.yaml; do
    cp "${f}" "${f/\.yaml/-${COMMON_INSTANCETYPES_VERSION}\.yaml}"
done
