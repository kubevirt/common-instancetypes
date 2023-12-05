# Development

This doc outlines how to build `common-instancetypes` and which dependencies are required.

## Requirements

A `Makefile` with common operations is supplied with this repository.

The following tools are required by each `Makefile` target:

* `make`
  * `make lint`
    * [yamllint](https://github.com/adrienverge/yamllint)
    * [shellcheck](https://github.com/koalaman/shellcheck)
  * `make generate`
    * [kustomize](https://kustomize.io/)
    * [yq](https://github.com/mikefarah/yq)
  * `make readme`
    * [yq](https://github.com/mikefarah/yq)
  * `make schema`
    * [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema)
  * `make validate`
    * [kustomize](https://kustomize.io/)
    * [kubeconform](https://github.com/yannh/kubeconform)
  * `make test-lint`
    * [golangci-lint](https://golangci-lint.run/)

By default `make` runs the `lint`, `validate`, `readme`, `test-lint` and `test` targets and their dependencies.

## Provided Makefile targets

The following targets are provided by the `Makefile`

* `make lint`: Run linting jobs on source and generated files.
* `make generate`: Generate the common-instancetypes bundles. Output can be found in `./_build`.
* `make validate`: Validate generated files against the OpenAPI schema of KubeVirt.
* `make schema`: Grab and convert the KubeVirt OpenAPI schema to a JSON schema to be used by `validate`.
* `make test`: Run unit tests against source files.
* `make test-fmt`: Run go fmt against test source files.
* `make test-vet`: Run go vet against test source files.
* `make test-lint`: Run golangci-lint against test source files.
* `make cluster-functest`: Run functests against latest stable KubeVirt.
* `make kubevirt-functest`: Run functests against latest `main` KubeVirt.
* `make readme`: Regenerate the list of provided resources at the end of `README.md`.
* `make clean`: Clean up generated files by other `Makefile` targets.
* `make build_image`: Build container image which provides build dependencies.
* `make push_image`: Build and push container image which provides build dependencies.
* `make kustomize`: Install vendored `kustomize` from `./tools` into `./_bin`.
* `make kubeconform`: Install vendored `kubeconform` from `./tools` into `./_bin`.
* `make yq`: Install vendored `yq` from `./tools` into `./_bin`.
* `make golangci-lint`: Install `golangci-lint` version specified in `Makefile` into `./_bin`.

## Running Makefile targets within a container

A [`Dockerfile`](./Dockerfile) is available that provides the `openapi2jsonschema`, `shellcheck` and `yamllint` dependencies. The golang based depedencies are vendored in the `./tools` directory. To run them golang is required, which is also provided by the [`Dockerfile`](./Dockerfile).

To build and use this simply set the `COMMON_INSTANCETYPES_CRI` env variable to `docker` or `podman` before running `make build_image`:

```
$ COMMON_INSTANCETYPES_CRI=podman make build_image
$ COMMON_INSTANCETYPES_CRI=podman make
```

## Building the common-instancetypes bundles

Run the following command to build the common-instancetypes bundles:

```
make
```

The resulting bundles can be found in the `./_build` directory.

## Installation

[kustomize](https://kustomize.io/) is used to generate these resources before applying them to the cluster with [`kubectl`](https://kubernetes.io/docs/reference/kubectl/):

```
$ kustomize build | kubectl apply -f -
```

A number of resource specific [kustomize](https://kustomize.io/) build configs are also provided and allow limited generation. For instance the following command will generate and install only the [`VirtualMachineInstancetype`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineinstancetype) resources:

```
$ kustomize build VirtualMachineInstancetypes | kubectl apply -f -
```

Additional targets are avilable for [`VirtualMachineClusterInstancetypes`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineclusterinstancetype), [`VirtualMachineClusterPreferences`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineclusterpreference) and [`VirtualMachinePreferences`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachinepreference).

## Spinning up a cluster for testing purposes

This repository contains following methods to spin up a cluster for testing purposes:

* `make cluster-{up|sync|down}`
  * Spin up a kubevirtci cluster and install the latest release of KubeVirt.
* `make kubevirt-{up|sync|down}`
  * Spin up a kubevirtci cluster and build and install the latest `main` KubeVirt.

The `sync` target syncs resources on the cluster with the local tree.

For more documentation on the usage of kubevirtci clusters see [kubevirtci's README.md](https://github.com/kubevirt/kubevirtci/blob/main/README.md).
