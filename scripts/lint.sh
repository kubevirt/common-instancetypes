#!/bin/bash

yamllint .
bashate -i E006 ./scripts/*
