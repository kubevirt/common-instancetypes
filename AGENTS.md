# AGENTS.md

Instructions for AI coding agents working in this repository.

## Overview & scope

This repo ships [KubeVirt](https://kubevirt.io) [VirtualMachineInstancetype](https://kubevirt.io/user-guide/virtual_machines/instancetypes/) and [VirtualMachinePreference](https://kubevirt.io/user-guide/virtual_machines/instancetypes/) manifests as kustomize-built YAML. Human-facing docs: [README.md](./README.md), [DEVELOPMENT.md](./DEVELOPMENT.md), [docs/labels.md](./docs/labels.md).

This file applies to the whole repo. If nested `AGENTS.md` files are added later, the closest file to the edited path wins; explicit user instructions override both.

| Path                                                          | Purpose                                                                           |
| ------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| `instancetypes/`                                              | Series definitions (u, o, cx, m, n, rt, d) and sizes                              |
| `preferences/`                                                | Guest OS / workload preferences (Linux, Windows, RHEL, etc.)                      |
| `preferences/components/`                                     | Reusable preference fragments (EFI, disk bus, …)                                  |
| `VirtualMachine*Instancetypes` / `VirtualMachine*Preferences` | Top-level kustomize entrypoints for partial installs                              |
| `scripts/`                                                    | Generate, lint, validate, cluster helpers                                         |
| `tests/`                                                      | Go unit tests (`unittests/`) and cluster functests (`functests/`)                 |
| `tools/`                                                      | Vendored Go modules used to install `kustomize`, `yq`, `kubeconform` into `_bin/` |
| `_build/`                                                     | **Generated** bundles (do not edit by hand)                                       |
| `_bin/`                                                       | **Generated** local tool binaries                                                 |

## Build, test & validation commands

Run from the repository root unless noted. Prefer `make` targets over inventing commands.

```bash
# Default CI-style check (lint, validate, readme, test-lint, unit tests)
make

# Regenerate bundles under ./_build after changing source YAML/kustomize
make generate

# Validate generated manifests against KubeVirt OpenAPI schema
make validate

# Refresh README resource tables (required when adding/removing resources)
make readme

# Fail if generated output or README is out of date
make check-tree-clean

# Scoped checks
make lint          # yamllint, shellcheck (via scripts/cri.sh)
make test          # cd tests && go test ./unittests/...
make test-lint     # golangci-lint + go vet on tests/
make test-fmt      # gofumpt on tests/
```

Containerized toolchain (optional):

```bash
COMMON_INSTANCETYPES_CRI=podman make build_image
COMMON_INSTANCETYPES_CRI=podman make
```

Cluster functests (need a running cluster; slow):

```bash
make cluster-up && make cluster-sync && make cluster-functest
# or against KubeVirt main: make kubevirt-up && make kubevirt-sync && make kubevirt-functest
```

Install vendored CLI tools into `./_bin` (added to `PATH` by the Makefile):

```bash
make kustomize yq kubeconform golangci-lint
```

## Conventions & patterns

### Instance type naming

Follow the schema in [README.md](./README.md): `{class}{version}.{size}` (e.g. `u1.xlarge`, `cx1.medium1gi`). Classes: `u`, `o`, `cx`, `m`, `n`, `rt`, `d` (see series table in README).

### Source layout

- Series live under `instancetypes/<class>/<version>/` (e.g. `instancetypes/u/1/`).
- Sizes are often in `sizes.yaml` or series-specific YAML; kustomize wires them together.
- Preferences mirror OS/vendor paths under `preferences/` (e.g. `preferences/rhel/9/amd64/`).
- Shared preference building blocks go in `preferences/components/`.

### Labels

When adding or changing instance types, set labels documented in [docs/labels.md](./docs/labels.md) (`instancetype.kubevirt.io/cpu`, `memory`, `size`, class-specific optional labels, etc.). Version labels on bundles are injected at generate time—do not duplicate that logic in hand-written source.

### API version

Resources use `apiVersion: instancetype.kubevirt.io/v1beta1` unless the repo has already migrated; match neighboring files in the same directory.

### Go tests

- Edit tests under `tests/unittests/` or `tests/functests/`.
- Do **not** edit `tests/vendor/` or `tools/vendor/`.
- Match existing Ginkgo/Gomega style in functests.
- Use `DescribeTable` or `Entry` for parameterized tests whenever possible; avoid duplicating test logic across individual `It` blocks.

### Scripts

Shell scripts live in `scripts/` and are invoked by the Makefile. Keep them executable; `shellcheck` runs in `make lint`.

## Adding a preference

Playbook for agents (and humans). When the user gives minimal input (e.g. “add preference `foo` like `debian`”), copy the named template, apply their deltas, then run validation below.

### Minimum input from the user

Ask only if not inferable from “like `debian`” / “like `windows.2k25`” or sibling preferences:

| Field             | Example                                                                   |
| ----------------- | ------------------------------------------------------------------------- |
| **Name**          | `debian`, `windows.2k25.virtio`, `rhel.9.arm64`                           |
| **Template**      | “like debian”, “like windows 2k25 virtio”, “like rhel.9 amd64”            |
| **Requirements**  | CPU/memory (or “same as rhel.9”)                                          |
| **Boot functest** | yes/no; container image; SSH user (Linux) or Windows checks               |
| **Architectures** | `amd64`, `arm64`, `s390x` (multi-arch → separate dirs or arch components) |

Skip boot functests unless the user wants them **and** a container disk exists in CI (often `registry:5000/...` for Windows/RHEL).

### Copy-from templates

| Kind                 | Copy from                                                  | Resulting preference name                             |
| -------------------- | ---------------------------------------------------------- | ----------------------------------------------------- |
| Generic Linux (EFI)  | `preferences/debian/`                                      | `metadata.name` patch → e.g. `debian`                 |
| Linux (no EFI)       | `preferences/linux/` subtree                               | e.g. `linux`, `legacy`                                |
| Virtio transitional  | `preferences/linux-virtio-transitional/`                   | `linux.virtiotransitional`                            |
| Multi-arch Linux     | `preferences/fedora/` or `preferences/rhel/9/`             | per-arch dirs + `preferences/<os>/kustomization.yaml` |
| Windows (base)       | `preferences/windows/2k25/`                                | `windows.2k25` via `nameSuffix` on `../base`          |
| Windows (virtio)     | `preferences/windows/2k25_virtio/`                         | `windows.2k25.virtio`                                 |
| Windows variant only | `preferences/windows/11_virtio_dedicated/`                 | extends sibling; often **no** new functest            |
| DPDK / realtime      | `preferences/rhel/9_dpdk/`, `preferences/rhel/9_realtime/` | extra requirement labels + components                 |

### Files to create or update

**Every new preference family** (typical Linux single-name layout):

```
preferences/<dir>/
  kustomization.yaml          # resources + components + name patches
  metadata/kustomization.yaml # Component patches (often identical to debian)
  metadata/metadata.yaml      # display-name, tags, iconClass
  requirements/kustomization.yaml
  requirements/requirements.yaml  # spec.requirements + required-cpu/memory labels
```

Also add `- ./<dir>` to the parent `resources:` list:

- Top-level OS: `preferences/kustomization.yaml`
- Windows only: `preferences/windows/kustomization.yaml`

**Boilerplate that rarely changes:** `metadata/kustomization.yaml` and `requirements/kustomization.yaml` match `preferences/debian/`; only `metadata.yaml` and `requirements.yaml` are OS-specific.

**New shared behavior** (unusual NIC, disk bus, PCI hole): add under `preferences/components/`, reference from the preference `kustomization.yaml` `components:` list (see `preferences/windows/xp/` → `pcihole64`).

**Requirements labels** must match `spec.requirements` (unit tests enforce this). See [docs/labels.md](./docs/labels.md) for `instancetype.kubevirt.io/required-cpu` and `required-memory`.

### Kustomize patterns (quick reference)

**Linux (single name, EFI):** base `../linux-efi`, patch `/metadata/name` to the final preference name on both `VirtualMachinePreference` and `VirtualMachineClusterPreference`.

**Windows:** base `../base` (under `preferences/windows/<version>/`), `nameSuffix: ".<version>"` produces `windows.<version>`. Virtio dir: `resources: [../<version>]`, virtio components, `nameSuffix: .virtio`.

**Multi-arch:** shared `requirements/` at family root; per-arch subdirs add `preferred-architecture` / `required-architecture` components and their own `metadata/` (see `preferences/fedora/amd64/`).

### What is already automatic (do not do by hand)

| Task                                          | How                                                                                                                                                              |
| --------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| README preference table                       | `make readme` only — never edit the table in README.md manually                                                                                                  |
| Unit tests for labels/requirements            | `tests/unittests/` loads `_build/common-clusterpreferences-bundle.yaml`                                                                                          |
| “All preferences with requirements” functests | Dynamic list in `tests/functests/instancetype_test.go`; add to `skipPreference` only for exceptions (`legacy`, `linux`, `linux.efi`, `linux.virtiotransitional`) |

### Boot functests (optional, manual wiring)

Only when a boot test is requested and a container disk is available. Edit **two** files:

1. `tests/functests/test_suite_test.go` — `default*ContainerDisk` const, package var, `flag.StringVar` in `init()`
2. `tests/functests/instancetype_test.go` — one `Entry(...)` in `DescribeTable("a Linux guest with"` or `"a Windows guest with"`

Entry shape: `containerDisk`, `instancetype` (e.g. `u1.small`, `u1.nano`), `map[string]string` arch → preference name, `[]testFn` (e.g. `expectSSHToRunCommandOnLinux("cloud-user")`, `expectGuestAgentToBeConnected`).

Match naming/style of the nearest existing `Entry` (Fedora, Debian, Windows 2k25, etc.). Custom images often use `registry:5000/...`; public Linux images often use `quay.io/containerdisks/...`.

Upstream sometimes splits preference and functest into two commits; either order is fine if `make` passes.

### Agent workflow (end-to-end)

1. Identify template from user message or closest existing preference.
2. Create/update kustomize tree and parent `kustomization.yaml` `resources:`.
3. Run from repo root:

```bash
make generate validate test readme
make test-lint   # if Go tests changed
```

4. If functests were added and a cluster is available: `make cluster-sync && make cluster-functest` (slow; user may run this).
5. Before PR: `make check-tree-clean`.

### When not to add functests

- Variant preferences (e.g. dedicated virtio) that reuse an already-tested guest.
- No container disk image in CI yet (preference YAML only; functest follow-up later).
- User explicitly asks for preference manifests only.

## Workflow for typical changes

1. Edit YAML/kustomize under `instancetypes/` and/or `preferences/`.
2. Run `make generate validate test`.
3. If resources were added, removed, or renamed: `make readme`.
4. Before opening a PR: `make check-tree-clean`.

## Pull requests

Use [.github/PULL_REQUEST_TEMPLATE.md](./.github/PULL_REQUEST_TEMPLATE.md): describe what/why, link issues, fill the `release-note` block (`NONE` if not needed).

## Safety & guardrails

- **Do not** commit secrets, kubeconfigs, or cluster credentials.
- **Do not** hand-edit files in `_build/`, `_bin/`, `_cluster-up/`, `_kubevirt/`, or `_schemas/`.
- **Do not** modify vendored trees (`tests/vendor/`, `tools/vendor/`) except via upstream vendor bumps.
- **Do not** run destructive cluster or git commands unless the user explicitly asks.
- **Do not** push or open PRs unless the user asks.
- Prefer minimal diffs: one series or one preference family per change when possible.
- When unsure about series characteristics (hugepages, dedicated CPU, DPDK node labels, etc.), read the series section in [README.md](./README.md) before changing specs.

## Where to look

| Task               | Start here                                                                                             |
| ------------------ | ------------------------------------------------------------------------------------------------------ |
| Add a VM size      | `instancetypes/<series>/`, then `make readme`                                                          |
| Add OS preference  | [Adding a preference](#adding-a-preference); copy `preferences/debian/` or `preferences/windows/2k25/` |
| Label semantics    | [docs/labels.md](./docs/labels.md)                                                                     |
| Makefile targets   | [DEVELOPMENT.md](./DEVELOPMENT.md), [Makefile](./Makefile)                                             |
| Install to cluster | `kustomize build \| kubectl apply -f -` (see DEVELOPMENT.md)                                           |
