if [ "${COMMON_INSTANCETYPES_CRI}" == "" ]; then
    echo "COMMON_INSTANCETYPES_CRI must be set to either docker or podman"
    exit 1
fi

${COMMON_INSTANCETYPES_CRI} push ${COMMON_INSTANCETYPES_IMAGE}:v$(date +%Y%m%d)-$(git rev-parse --short HEAD)