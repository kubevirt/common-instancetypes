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

FROM quay.io/fedora/fedora-minimal:39

RUN microdnf install -y make golang python3-pip ShellCheck && microdnf clean all -y

# yamllint openapi2jsonschema
RUN pip install yamllint && \
  # Manually provide deps of openapi2jsonschema to force it to use a later pyyaml version
  pip install --no-dependencies openapi2jsonschema jsonref click && \
  pip cache purge
