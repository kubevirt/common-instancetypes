#!/usr/bin/env bash

set -e

if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    /bin/bash -c "$@"
    exit $?
fi

${COMMON_INSTANCETYPES_CRI} run -v "$(pwd):/common-instancetypes:rw,Z" -e KUBEVIRT_VERSION="${KUBEVIRT_VERSION}" -e KUBEVIRT_TAG="${KUBEVIRT_TAG}" --rm "${COMMON_INSTANCETYPES_IMAGE}:${COMMON_INSTANCETYPES_IMAGE_TAG}" /usr/bin/bash -c "cd /common-instancetypes && $*"

