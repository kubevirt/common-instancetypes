#!/bin/bash

_base_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
context_dir="${_base_dir}/image"

curl -L "https://yum.oracle.com/templates/OracleLinux/OL8/u10/x86_64/OL8U10_x86_64-kvm-b258.qcow2" -o "${context_dir}/OL8U10-amd64.qcow2"
curl -L "https://yum.oracle.com/templates/OracleLinux/OL8/u10/aarch64/OL8U10_aarch64-kvm-cloud-b123.qcow2" -o "${context_dir}/OL8U10-arm64.qcow2"
curl -L "https://yum.oracle.com/templates/OracleLinux/OL9/u5/x86_64/OL9U5_x86_64-kvm-b259.qcow2" -o "${context_dir}/OL9U5-amd64.qcow2"
curl -L "https://yum.oracle.com/templates/OracleLinux/OL9/u5/aarch64/OL9U5_aarch64-kvm-cloud-b127.qcow2" -o "${context_dir}/OL9U5-arm64.qcow2"

declare -A version_images=(
    ["8.10"]="OL8U10-amd64.qcow2 OL8U10-arm64.qcow2"
    ["9.5"]="OL9U5-amd64.qcow2 OL9U5-arm64.qcow2"
  )


for version in "${!version_images[@]}"; do
  IFS=' ' read -r -a arch_images <<< "${version_images[$version]}"
  manifest_name="oraclelinux:${version}"


  echo "Processing version: ${version}"

  # Build per-arch images
  for image in "${arch_images[@]}"; do
    arch="${image##*-}"        # Extracts "amd64" or "arm64" (removes .qcow2)
    arch="${arch%.*}"          # Strips .qcow2
    temp_tag="${manifest_name}-${arch}"

    echo "Building ${arch} image with disk: ${image}"
    podman build \
      --platform="linux/${arch}" \
      --format "docker" \
      --build-arg disk="${image}" \
      -t "${temp_tag}" \
      -f "${context_dir}/oraclelinux.Dockerfile" \
      "${context_dir}"

 done

  echo "Creating multi-arch manifest: ${manifest_name}"
  podman manifest create "${manifest_name}"
  podman manifest add "${manifest_name}" "${manifest_name}-amd64"
  podman manifest add "${manifest_name}" "${manifest_name}-arm64"

  echo "Pushing multi-arch manifest: ${manifest_name}"
done
