# A basic container image for the following Makefile target deps:
#
#* `make`
#  * `make lint`
#    * [yamllint](https://github.com/adrienverge/yamllint)
#    * [ShellCheck](https://github.com/koalaman/shellcheck)
#  * `make generate`
#    * [kustomize](https://kustomize.io/)
#  * `make validate`
#    * [kustomize](https://kustomize.io/)
#    * [kubeconform](https://github.com/yannh/kubeconform)
#  * `make readme`
#    * [yq](https://github.com/mikefarah/yq)
#* `make schema` (optional)
#  * [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema)

FROM quay.io/fedora/fedora-minimal:39

RUN microdnf install -y make git golang python3-pip tar ShellCheck && microdnf clean all -y

# yamllint bashate openapi2jsonschema
RUN pip install yamllint bashate && \
  # Manually provide deps of openapi2jsonschema to force it to use a later pyyaml version
  pip install --no-dependencies openapi2jsonschema jsonref click && \
  pip cache purge

# kustomize and kubeconform
RUN go install sigs.k8s.io/kustomize/kustomize/v5@v5.2.1 && \
  go install github.com/yannh/kubeconform/cmd/kubeconform@v0.6.3 && \
  mv /root/go/bin/kustomize /usr/local/bin/kustomize && \
  mv /root/go/bin/kubeconform /usr/local/bin/kubeconform && \
  rm -rf /root/go

# Download latest yq binary
RUN curl -L https://github.com/mikefarah/yq/releases/download/v4.35.2/yq_linux_amd64 -o /usr/local/bin/yq && \
  chmod +x /usr/local/bin/yq
