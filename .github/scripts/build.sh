#!/bin/bash

if [ -z "${COMMON_INSTANCETYPES_VERSION}" ] ; then
    echo "${BASH_SOURCE[0]} expects the following env variables to be provided: COMMON_INSTANCETYPES_VERSION."
    exit 1
fi

# Build the bundles
make && COMMON_INSTANCETYPES_VERSION=${COMMON_INSTANCETYPES_VERSION} make generate

for f in *bundle.yaml; do
    cp "${f}" "${f/\.yaml/-${COMMON_INSTANCETYPES_VERSION}\.yaml}"
done
