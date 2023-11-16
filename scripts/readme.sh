#!/bin/bash

set -e

if ! command -v yq &> /dev/null; then
    echo "yq is not installed, see https://github.com/mikefarah/yq#install for more details."
    exit 1
fi

sed -i '/\#\# Resources/Q' README.md
resources=$(yq '.kind, .metadata.name' _build/common-instancetypes-all-bundle.yaml)
cat <<EOF >> README.md
## Resources

The following instance type and preference resources are provided by this project:

\`\`\`

$resources

\`\`\`
EOF
