#!/bin/bash

if !command -v yamllint &> /dev/null; then
    echo "yamllint is not installed, see https://github.com/adrienverge/yamllint#installation for more details."
    exit 1
fi

if !command -v bashate &> /dev/null; then
    echo "bashate is not installed, see https://docs.openstack.org/bashate/latest/install/index.html for more details."
    exit 1
fi

yamllint .

bashate -i E006 ./scripts/*
