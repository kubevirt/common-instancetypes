#!/usr/bin/env bash

if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    /bin/bash -c "$@"
    exit $?
fi

${COMMON_INSTANCETYPES_CRI} run -v "$(pwd):/common-instancetypes" -e KUBEVIRT_VERSION=${KUBEVIRT_VERSION} --rm ${COMMON_INSTANCETYPES_IMAGE} /usr/bin/bash -c "cd /common-instancetypes &&  $@"
