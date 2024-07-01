#!/usr/bin/env bash

set -e

# Default to podman if COMMON_INSTANCETYPES_CRI isn't provided
COMMON_INSTANCETYPES_CRI=${COMMON_INSTANCETYPES_CRI:-podman}

if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    /bin/bash -c "$@"
    exit $?
fi

${COMMON_INSTANCETYPES_CRI} run -v "$(pwd):/common-instancetypes:rw,Z" -e KUBEVIRT_VERSION="${KUBEVIRT_VERSION}" --rm "${COMMON_INSTANCETYPES_IMAGE}:${COMMON_INSTANCETYPES_IMAGE_TAG}" /usr/bin/bash -c "cd /common-instancetypes && $*"

