# A basic container image for the following Makefile target deps:
#
#* `make`
#  * `make lint`
#    * [yamllint](https://github.com/adrienverge/yamllint)
#    * [bashate](https://github.com/openstack/bashate)
#  * `make generate`
#    * [kustomize](https://kustomize.io/)
#  * `make validate`
#    * [kustomize](https://kustomize.io/)
#    * [kubeconform](https://github.com/yannh/kubeconform)
#  * `make readme`
#    * [yq](https://github.com/mikefarah/yq)
#* `make schema` (optional)
#  * [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema)

FROM quay.io/fedora/fedora-minimal:36

RUN microdnf install -y make git golang python3-pip tar && microdnf clean all -y

# yamllint bashate openapi2jsonschema
RUN pip install yamllint bashate openapi2jsonschema && pip cache purge

# kustomize
RUN curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/kustomize%2Fv4.5.7/hack/install_kustomize.sh"  | bash && mv /kustomize /bin/kustomize

# kubeconform
RUN go install github.com/yannh/kubeconform/cmd/kubeconform@v0.5.0 && mv /root/go/bin/kubeconform /bin/kubeconform

# Download latest yq binary
RUN curl -L https://github.com/mikefarah/yq/releases/download/v4.29.2/yq_linux_amd64 -o /usr/local/bin/yq && \
    chmod +x /usr/local/bin/yq

# We need to create a directory for the code and add it as a safe.directory before use, allowing git commands to work correctly.
RUN mkdir /common-instancetypes && git config --global --add safe.directory /common-instancetypes
