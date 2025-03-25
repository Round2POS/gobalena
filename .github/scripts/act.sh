#!/bin/bash

# Run this script to run (some of the) GitHub Actions locally.

# https://gist.github.com/mohanpedala/1e2ff5661761d3abd0385e8223e16425
set -e -x -v -u -o pipefail

PROJ_PATH=$(pwd)
TMP_DIR=$(mktemp -d)
function cleanup {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

ACTION_CACHE_PATH="${PWD}/.cache/act/action-cache"
CACHE_SERVER_PATH="${PWD}/.cache/act/cache-server-path"
# ACT_PROJECT_PATH="${PWD}/.cache/act/project-clone-path"
ACT_PROJECT_PATH="${TMP_DIR}/project-clone-path"

################################################################################
mkdir -p "${ACT_PROJECT_PATH}"
mkdir -p "${ACTION_CACHE_PATH}"
mkdir -p "${CACHE_SERVER_PATH}"

chmod -R 777 "${ACT_PROJECT_PATH}" || true
rm -Rf "${ACT_PROJECT_PATH}" || true
mkdir -p "${ACT_PROJECT_PATH}"

################################################################################
# Set aside a directory to use as the repo, for the GH action.

# Save the staged changes to a patch file.
git diff --cached --binary > "${TMP_DIR}/staged_changes.patch"
# Make a clean copy of the project directory.
# Note, for a local repo, this retains the current branch and commit, but not
# the staged changes.
git clone "${PROJ_PATH}" "${ACT_PROJECT_PATH}"

# An alternative method, but the resulting directory won't be a git repo.
# git checkout-index --all --prefix="${ACT_PROJECT_PATH}/"
# chmod -R 555 "${ACT_PROJECT_PATH}"

cd "${ACT_PROJECT_PATH}"
# Apply the staged changes to the clean copy.
git apply --whitespace=fix --allow-empty "${TMP_DIR}/staged_changes.patch"
# Stage the changes.
git add .
################################################################################

# Use --bind to keep the directory persistent, useful if you need to check the
# contents of the directory after the run for why it failed. However, if you
# do use --bind, you will need to manually clean up the directory after the run,
# because it uses root to create the directory and files.
go run github.com/nektos/act@v0.2.68 \
  -P "ubuntu-v22.04-c16-m64=catthehacker/ubuntu:act-22.04" \
  push \
  --action-offline-mode \
  --use-new-action-cache \
  --cache-server-path "${CACHE_SERVER_PATH}" \
  --action-cache-path "${ACTION_CACHE_PATH}" \
  --job build-and-test
################################################################################
