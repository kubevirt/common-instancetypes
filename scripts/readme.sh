#!/bin/bash
sed -i '/\#\# Resources/Q' README.md
resources=$(yq '.kind, .metadata.name' common-instancetypes-all-bundle.yaml)
cat <<EOF >> README.md
## Resources

The following instance type and preference resources are provided by this project:

\`\`\`

$resources

\`\`\`
EOF