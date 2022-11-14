#!/bin/bash

if ! command -v yamllint &> /dev/null; then
    echo "yamllint is not installed, see https://github.com/adrienverge/yamllint#installation for more details."
    exit 1
fi

if ! command -v shellcheck &> /dev/null; then
    echo "shellcheck is not installed, see https://github.com/koalaman/shellcheck#installing for more details."
    exit 1
fi

yamllint .

shellcheck ./scripts/*
