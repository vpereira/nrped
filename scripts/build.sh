#!/bin/bash
#
# This script builds the application from source.
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd $DIR

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# If we're building on Windows, specify an extension
EXTENSION=""
if [ "$(go env GOOS)" = "windows" ]; then
    EXTENSION=".exe"
fi

# Build!
echo "--> Building..."
cd common && go build common.go && cd ..
cd read_config && go build read_config.go && cd ..
cd check_nrpe && go build -o check_nrpe && cd ..
cd gen_certificate && go build -o gen_certificate && cd ..
go build  -o nrped
#cp nrped $GOPATH/bin
#cp check_nrpe/check_nrpe $GOPATH/bin
#go build \
#    -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}" \
#   -v \
#    -o bin/nrpe${EXTENSION}
#cp bin/nrpe${EXTENSION} $GOPATH/bin
