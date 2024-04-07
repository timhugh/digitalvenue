#!/usr/bin/env sh

root=$(git rev-parse --show-toplevel)

ln -s "${root}/.scripts/pre-commit" "${root}/.git/hooks/pre-commit"
