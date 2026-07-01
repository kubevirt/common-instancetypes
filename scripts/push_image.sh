#!/bin/bash

set -e

if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    echo "COMMON_INSTANCETYPES_CRI must be set to either docker or podman"
    exit 1
fi

tag="v$(date +%Y%m%d)-$(git rev-parse --short HEAD)"

${COMMON_INSTANCETYPES_CRI} build image -t "${COMMON_INSTANCETYPES_IMAGE}" -t "${COMMON_INSTANCETYPES_IMAGE}:${tag}"

${COMMON_INSTANCETYPES_CRI} push "${COMMON_INSTANCETYPES_IMAGE}"
${COMMON_INSTANCETYPES_CRI} push "${COMMON_INSTANCETYPES_IMAGE}:${tag}"
