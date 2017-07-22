#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
pushd $GOPATH/src/github.com/appscode/searchlight

if [ -n "$1" ]; then
	echo "Using BUILD_BRANCH=release-$1"
	echo ""
    git checkout release-$1
    ./hack/make.py build
    APPSCODE_ENV=prod ./hack/make.py push
    ./hack/make.py push

    ./hack/docker/searchlight/setup.sh
    APPSCODE_ENV=prod ./hack/docker/searchlight/setup.sh release

    ./hack/docker/icinga/build.sh
    ./hack/docker/icinga/build.sh release
    
    ./hack/docker/icinga/setup.sh
    ./hack/docker/icinga/setup.sh release
else
	echo "No release tag specified."
fi
