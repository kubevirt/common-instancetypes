#!/bin/bash

if !command -v openapi2jsonschema &> /dev/null; then
    echo "openapi2jsonschema is not installed, see https://github.com/mikefarah/yq#install for more details."
    exit 1
fi

KUBEVIRT_SWAGGER_URL=https://raw.githubusercontent.com/kubevirt/kubevirt/${KUBEVIRT_VERSION}/api/openapi-spec/swagger.json
SCHEMA_DIR=schemas/${KUBEVIRT_VERSION}
mkdir -p ${SCHEMA_DIR}

openapi2jsonschema --stand-alone --expanded --strict -o ${SCHEMA_DIR} ${KUBEVIRT_SWAGGER_URL}
