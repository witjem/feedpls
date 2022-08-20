#!/usr/bin/env bash
ROOT_DIR=$(git rev-parse --show-toplevel)
GIT_DIR=${ROOT_DIR}/.git
GIT_HOOKS_DIR=${GIT_DIR}/hooks

rm -rf "${GIT_HOOKS_DIR}"

ln -s ${ROOT_DIR}/githooks ${GIT_HOOKS_DIR}
