#!/bin/bash

if [ -z "${COMMON_INSTANCETYPES_VERSION}" ] || [ -z "${COMMON_INSTANCETYPES_RELEASE_URL}" ] || [ -z "${RELEASE_FORK_USER}" ] || [ -z "${GH_TOKEN}" ]; then
    echo "${BASH_SOURCE[0]} expects the following env variables to be provided: COMMON_INSTANCETYPES_VERSION, COMMON_INSTANCETYPES_RELEASE_URL, RELEASE_FORK_USER and GH_TOKEN."
    exit 1
fi

# Set git configs to sign the commit
git config --global user.email "${RELEASE_FORK_USER}@redhat.com"
git config --global user.name "common-instancetypes Release Automation"

# Clone repo and use a token to allow pushing before creating a PR
git clone https://github.com/kubevirt/ssp-operator && cd ssp-operator || exit
git checkout origin/master -b update-common-instancetypes-"${COMMON_INSTANCETYPES_VERSION}"

# Update the new common-instancetypes file
cp ../common-clusterinstancetypes-bundle-"${COMMON_INSTANCETYPES_VERSION}".yaml data/common-instancetypes-bundle/
cp ../common-clusterpreferences-bundle-"${COMMON_INSTANCETYPES_VERSION}".yaml data/common-instancetypes-bundle/

git add data && git commit -sm "common-instancetypes: Update bundle to COMMON_INSTANCETYPES_VERSION ${COMMON_INSTANCETYPES_VERSION}"

git remote add release https://"${RELEASE_FORK_USER}":"${GH_TOKEN}"@github.com/"${RELEASE_FORK_USER}"/ssp-operator
git push -f -u release update-common-instancetypes-"${COMMON_INSTANCETYPES_VERSION}"

# Create a new PR in the operator repo
gh pr create --repo kubevirt/ssp-operator \
    --base master \
    --head "${RELEASE_FORK_USER}":update-common-instancetypes-"${COMMON_INSTANCETYPES_VERSION}" \
    --title "Update common-instancetypes to ${COMMON_INSTANCETYPES_VERSION}" \
    --body "$(cat << EOF
common-instancetypes: Update bundle to ${COMMON_INSTANCETYPES_VERSION}

${COMMON_INSTANCETYPES_RELEASE_URL}

**Release note**:
\`\`\`release-note
Update common-instancetypes bundle to ${COMMON_INSTANCETYPES_VERSION}
\`\`\`
EOF
)
"
