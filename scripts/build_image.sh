#!/bin/bash

set -e

if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    echo "COMMON_INSTANCETYPES_CRI must be set to either docker or podman"
    exit 1
fi

${COMMON_INSTANCETYPES_CRI} build image -t "${COMMON_INSTANCETYPES_IMAGE}"
