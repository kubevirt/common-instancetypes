#!/bin/bash

set -e

if ! command -v openapi2jsonschema &> /dev/null; then
    echo "openapi2jsonschema is not installed, please run 'pip install --no-dependencies openapi2jsonschema jsonref click'."
    exit 1
fi

KUBEVIRT_SWAGGER_URL=https://raw.githubusercontent.com/kubevirt/kubevirt/${KUBEVIRT_TAG}/api/openapi-spec/swagger.json
SCHEMA_DIR=_schemas/${KUBEVIRT_TAG}
mkdir -p "${SCHEMA_DIR}"

openapi2jsonschema --stand-alone --expanded --strict -o "${SCHEMA_DIR}" "${KUBEVIRT_SWAGGER_URL}"
