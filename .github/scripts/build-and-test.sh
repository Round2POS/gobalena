#!/bin/bash
# https://gist.github.com/mohanpedala/1e2ff5661761d3abd0385e8223e16425
set -e -x -v -u -o pipefail

(
RED=$'\e[0;31m'
GREEN=$'\e[0;32m'
BLUE=$'\e[0;34m'
NC=$'\e[0m'

PROJ_PATH=$(pwd)
TMP_DIR=$(mktemp -d)
function cleanup {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

set +u
source ~/.bashrc
source ~/.profile
set -u
################################################################################
go build ./...
################################################################################
echo -e "${GREEN}Success: all tests have passed.${NC}"
################################################################################
)