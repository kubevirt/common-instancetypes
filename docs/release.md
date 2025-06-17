# Release Process

The common-instancetypes project releases just ahead of the core KubeVirt Code Freeze to allow the resulting manifests to be pulled into the kubevirt/kubevirt project and virt-operator component.

## Release tagging

The initial release is always tagged as `v1.${y}.0` on main with a `release-1.${y}` branch then cut from this tag.

The following release workflow is provided to then generate the release notes and release artifacts:

<https://github.com/kubevirt/common-instancetypes/blob/67e413352b081deaec3ec504912947836532a731/.github/workflows/release.yaml>

## Sync with kubevirt/kubevirt

Releases of common-instancetypes are sync'd into the core kubevirt/kubevirt codebase by the following periodic automation:

<https://github.com/kubevirt/project-infra/blob/3c7dc66b31fa5b54b25c3b7712d1e1bd9fdc42da/github/ci/prow-deploy/files/jobs/kubevirt/common-instancetypes/common-instancetypes-periodics.yaml>

This automation must be updated when new releases of both common-instancetypes and kubevirt are cut.
