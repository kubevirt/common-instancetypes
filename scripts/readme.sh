#!/bin/bash

set -e

if ! command -v yq &> /dev/null; then
    echo "yq is not installed, see https://github.com/mikefarah/yq#install for more details."
    exit 1
fi

sed -i '/\#\# Resources/Q' README.md
instancetypes=$(\
  yq ea '[.]'  _build/common-instancetypes-bundle.yaml | \
  yq '.[] | [.metadata.name, .spec.cpu.guest, .spec.memory.guest] | join("  |  ")'
)
preferences=$(\
  yq ea '[.]'  _build/common-preferences-bundle.yaml | \
  yq '.[] | [.metadata.name, .metadata.annotations."openshift.io/display-name"] | join(" | ")'
)
cat <<EOF >> README.md
## Resources

The following instancetype resources (cluster-wide and namespaced) are
provided by this project:

Name | vCPUs | Memory
-----|-------|-------
$instancetypes

The following preference resources (cluster-wide and namespaced) are
provided by this project:

Name | Guest OS
-----|---------
$preferences
EOF
