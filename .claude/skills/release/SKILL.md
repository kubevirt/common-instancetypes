---
name: release
description: Guide through the common-instancetypes release process
user_invocable: true
arguments:
  - name: version
    description: "Release version (e.g. v1.7.0)"
    required: true
---

# Common Instancetypes Release

Guide the user through the release process for common-instancetypes version `$ARGUMENTS.version`.

## Instructions

1. Read `docs/release.md` for the full checklist and context.
2. Parse the version from `$ARGUMENTS.version` to determine:
   - `Y` — the minor version (e.g. `7` for `v1.7.0`)
   - Release branch: `release-1.${Y}`
   - Previous release branch: `release-1.${Y-1}`
   - Corresponding kubevirt release branch for the periodic: `release-1.${Y+1}`
3. Run pre-flight checks:
   - Compare `GIMME_GO_VERSION` in `image/Dockerfile` with the Go version in `tools/go.mod` — warn if they're incompatible.
   - Check that `GOLANGCI_LINT_VERSION` in the `Makefile` is recent enough to support the project's Go version.
   - Check that `COMMON_INSTANCETYPES_IMAGE_TAG` in the `Makefile` references an image built with the correct Go version.
4. Check current state:
   - Does `release-1.${Y}` branch already exist? (`git branch -a | grep release-1.${Y}`)
   - Does `v1.${Y}.0` tag already exist? (`git tag -l v1.${Y}.0`)
   - Does a `v1.${Y}.0` milestone exist? (`gh api repos/kubevirt/common-instancetypes/milestones`)
   - Does a `v1.${Y+1}.0` milestone exist?
   - Check project-infra for existing release jobs and periodics.
5. Walk through each checklist item interactively, confirming with the user before taking actions that push to remote repositories or create PRs.
6. For the project-infra steps (presubmit jobs, periodics), the files live in:
   - `../project-infra/github/ci/prow-deploy/files/jobs/kubevirt/common-instancetypes/`
   - Copy `common-instancetypes-presubmits.yaml` to `common-instancetypes-1.${Y}.yaml`, updating branch references to `release-1.${Y}` and appending `-1.${Y}` to all job names.
   - Update `COMMON_INSTANCETYPES_BRANCH` in the main `periodic-update-common-instancetypes-bundles` to `release-1.${Y}`.
   - Add `periodic-update-common-instancetypes-bundles-1.${Y-1}` to sync `release-1.${Y-1}` with kubevirt `release-1.${Y+1}` if that kubevirt branch exists.
7. After completing steps, update the milestone description with the checklist and PR references.
8. Use `origin-ssh` (not `origin`) for push operations as HTTPS auth is not configured.

## Recovery

If the release workflow fails:
1. Diagnose the failure from the GitHub Actions logs.
2. Fix on `main`, cherry-pick to `release-1.${Y}`.
3. Delete and re-push the tag: `git push origin-ssh :refs/tags/v1.${Y}.0 && git tag -d v1.${Y}.0 && git tag v1.${Y}.0 && git push origin-ssh v1.${Y}.0`
