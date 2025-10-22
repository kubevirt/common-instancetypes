# Release Process

The common-instancetypes project releases after the `beta.0` release of `kubevirt/kubevirt`:

<https://github.com/kubevirt/sig-release/blob/main/releases/release-checklist.md#kubevirt-pre-release-checklist>

## Checklist

* [ ] Cut `release-1.${Y}` branch from `main`
* [ ] Tag `v1.${Y}.0` release on `release-1.${Y}` branch and generate release
* [ ] Ensure the [release
workflow](https://github.com/kubevirt/common-instancetypes/blob/67e413352b081deaec3ec504912947836532a731/.github/workflows/release.yaml)
runs and succeeds
* [ ] Copy and rename main pre-submit jobs definitions to a new release specific
file
* [ ] Update `periodic-update-common-instancetypes-bundles` periodic to sync
latest release with `kubevirt/kubevirt` `main`
* [ ] Once `periodic-update-common-instancetypes-bundles` runs work to merge
new release into `kubevirt/kubevirt` `main`
* [ ] Introduce `periodic-update-common-instancetypes-bundles-1.${Y}` once a
new version of `kubevirt/kubevirt` is cut to keep the corresponding `release`
branch sync'd with our `releease-1.${Y}` branch
* [ ] Create a `v1.{Y+1}.0` milestone for the new release cycle and copy this
checklist into description
* [ ] Update milestone once release schedule for the corresponding
`kubevirt/kubevirt` release agreed on
