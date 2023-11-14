#!/usr/bin/env bash

set -e

if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    /bin/bash -c "$@"
    exit $?
fi

${COMMON_INSTANCETYPES_CRI} run -v "$(pwd):/common-instancetypes:rw,Z" -e KUBEVIRT_VERSION="${KUBEVIRT_VERSION}" --rm "${COMMON_INSTANCETYPES_IMAGE}" /usr/bin/bash -c "cd /common-instancetypes &&  $*"

