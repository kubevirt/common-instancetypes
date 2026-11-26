# Release Process

The common-instancetypes project releases before the code freeze of `kubevirt/kubevirt`:

<https://github.com/kubevirt/sig-release/blob/main/releases/release-checklist.md#kubevirt-pre-release-checklist>

## Pre-flight checks

Before starting the release, verify:

* The Go version in `image/Dockerfile` (`GIMME_GO_VERSION`) is compatible
  with `tools/go.mod`. If not, bump it, rebuild the builder image, and update
  `COMMON_INSTANCETYPES_IMAGE_TAG` in the `Makefile`.
* The `golangci-lint` version in the `Makefile` (`GOLANGCI_LINT_VERSION`)
  supports the Go version used by the project.

## Checklist

* [ ] Cut `release-1.${Y}` branch from `main`
* [ ] Record the corresponding KubeVirt minor release and verify the latest
published stable patch for that minor in the [KubeVirt
releases](https://github.com/kubevirt/kubevirt/releases). On the new branch,
set `KUBEVIRT_VERSION` to that fixed `vX.Y.Z` tag and `KUBEVIRT_TAG` to the
matching `release-X.Y` branch in both the `Makefile` and the direct defaults in
`scripts/kubevirtci.sh` and `scripts/kubevirt.sh`, respectively. Resolve the
default `KUBEVIRTCI_TAG` from that fixed KubeVirt tag's
`kubevirtci/cluster-up/version.txt` and verify the kubevirtci tag exists.
If the target KubeVirt release branch or stable tag is not yet available,
record an explicit follow-up in the release milestone to set and verify these
pins as soon as they are published.
* [ ] Verify release defaults and explicit overrides through both `make` and
direct script calls for `KUBEVIRT_VERSION`, `KUBEVIRT_TAG`, and `KUBEVIRTCI_TAG`.
Use fresh checkouts or inspect the selected refs: existing `_kubevirt` and
`_cluster-up` directories are reused by the helpers. Run
`make generate lint validate test test-lint` with the release defaults and run
cluster functional tests in CI. Track any checks awaiting the target KubeVirt
release alongside the pinning follow-up.
* [ ] Tag `v1.${Y}.0` release on `release-1.${Y}` branch and generate release
* [ ] Ensure the [release
workflow](https://github.com/kubevirt/common-instancetypes/blob/67e413352b081deaec3ec504912947836532a731/.github/workflows/release.yaml)
runs and succeeds
* [ ] Copy and rename main pre-submit job definitions to a new release specific
file in
[`kubevirt/project-infra`](https://github.com/kubevirt/project-infra/tree/main/github/ci/prow-deploy/files/jobs/kubevirt/common-instancetypes)
* [ ] Update `periodic-update-common-instancetypes-bundles` in
[`common-instancetypes-periodics.yaml`](https://github.com/kubevirt/project-infra/blob/main/github/ci/prow-deploy/files/jobs/kubevirt/common-instancetypes/common-instancetypes-periodics.yaml)
to sync latest release with `kubevirt/kubevirt` `main`
* [ ] Once `periodic-update-common-instancetypes-bundles` runs work to merge
new release into `kubevirt/kubevirt` `main`
* [ ] Introduce `periodic-update-common-instancetypes-bundles-1.${Y-1}` once a
new version of `kubevirt/kubevirt` is cut to keep the corresponding `release`
branch sync'd with our `release-1.${Y-1}` branch
* [ ] Create a `v1.${Y+1}.0` milestone for the new release cycle and copy this
checklist into description
* [ ] Update milestone once release schedule for the corresponding
`kubevirt/kubevirt` release agreed on

## Recovery

If the release workflow fails after tagging:

1. Fix the issue on `main` and cherry-pick to `release-1.${Y}`
2. Delete the remote tag: `git push origin :refs/tags/v1.${Y}.0`
3. Re-tag on the updated branch and push: `git tag v1.${Y}.0 && git push origin v1.${Y}.0`
