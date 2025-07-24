# A basic container image for the following Makefile target deps:
#
#* `make`
#  * `make lint`
#    * [yamllint](https://github.com/adrienverge/yamllint)
#    * [ShellCheck](https://github.com/koalaman/shellcheck)
#  * `make schema`
#    * [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema)
#
# Golang is provided to allow running other Makefile targets
FROM quay.io/kubevirtci/golang:v20250701-f32dbda
ENV GIMME_GO_VERSION=1.24.5

RUN dnf install -y make tar python3-pip ShellCheck && dnf clean all -y

# yamllint openapi2jsonschema
RUN pip install yamllint && \
  # Manually provide deps of openapi2jsonschema to force it to use a later pyyaml version
  pip install --no-dependencies openapi2jsonschema jsonref click && \
  pip cache purge
