# kubevirt/common-instancetypes

A set of instancetypes and preferences to create KubeVirt VMs.

## Requirements

- kustomize
- yamllint

## Build

```
$ make
yamllint .
echo "---" > common-instancetypes.yaml
kustomize build ./instancetypes >> common-instancetypes.yaml
echo "---" >> common-instancetypes.yaml
kustomize build ./preferences >> common-instancetypes.yaml
```
