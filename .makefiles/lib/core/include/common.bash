function log() {
    >&2 echo "---" "$@"
}

function die() {
    >&2 echo "!!! $@"
    exit 1
}

if [ -n "${MF_DEBUG:-}" ]; then
    set -x
fi

if [ -z "${MF_PROJECT_ROOT:-}" ]; then
    MF_PROJECT_ROOT="$(cd "$(dirname ${BASH_SOURCE[0]})/../../../../"; pwd)"
fi

if [ -z "${MF_ROOT:-}" ]; then
    MF_ROOT="$MF_PROJECT_ROOT/.makefiles"
fi

MF_VERSION="1"

if [[ "$(which build-resource)" != "$MF_ROOT/lib/core/bin/build-resource" ]]; then
    PATH="$MF_ROOT/lib/core/bin:$PATH"
fi
