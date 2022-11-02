#!/bin/bash

if [ "${COMMON_INSTANCETYPES_CRI}" != "" ]; then
    ${COMMON_INSTANCETYPES_CRI} build . -t ${COMMON_INSTANCETYPES_IMAGE}
fi