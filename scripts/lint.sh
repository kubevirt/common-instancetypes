#!/bin/bash

if ! command -v yamllint &> /dev/null; then
    echo "yamllint is not installed, see https://github.com/adrienverge/yamllint#installation for more details."
    exit 1
fi

if ! command -v shellcheck &> /dev/null; then
    echo "shellcheck is not installed, see https://github.com/koalaman/shellcheck#installing for more details."
    exit 1
fi

if ! yamllint -d relaxed instancetypes preferences common-*-bundle.yaml; then
    exit 1
fi

if ! shellcheck scripts/* .github/scripts/* ; then
    exit 1
fi
