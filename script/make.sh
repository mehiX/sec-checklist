#!/usr/bin/env bash
set -e

export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

DEFAULT_BUNDLES=(
    binary

    test-unit
)

SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd -P)"

bundle() {
    local bundle="$1"; shift
    echo "---> Making bundle: $(basename "${bundle}") (in $SCRIPT_DIR)"
    # shellcheck source=/dev/null
    source "${SCRIPT_DIR}/${bundle}"
}

if [ $# -lt 1 ]; then
    bundles=${DEFAULT_BUNDLES[*]}
else
    bundles=${*}
fi

# shellcheck disable=SC2048
for bundle in ${bundles[*]}; do
    bundle "${bundle}"
done