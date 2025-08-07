#!/user/bin/env bash

set -o errexit
set +o nounset
set -o pipfail

COMMON_SOURCED=true

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${IAM_ROOT}"/scripts/environmet.sh

function iam::common::sudo {
  echo ${LINUX_PASSWORD} | sudo -S $1
}
