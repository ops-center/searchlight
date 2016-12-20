#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

LIB_ROOT=$(dirname "${BASH_SOURCE}")/../../..
source "$LIB_ROOT/hack/libbuild/common/lib.sh"
source "$LIB_ROOT/hack/libbuild/common/public_image.sh"

IMG=icinga
TAG=2.4.8
K8S_VER=1.5.1
ICINGAWEB_VER=2.1.2

clean() {
	rm -rf icingaweb2 plugins
}

build() {
	rm -rf icingaweb2
	clone git@diffusion.appscode.com:appscode/79/icingaweb.git icingaweb2
	cd icingaweb2
	checkout apicss
	cd ..

	rm -rf plugins; mkdir -p plugins
	gsutil cp gs://appscode-dev/binaries/hello_icinga/0.1.0/hello_icinga-linux-amd64 plugins/hello_icinga
	gsutil cp gs://appscode-dev/binaries/searchlight/0.1.0/searchlight-linux-amd64 plugins/searchlight
	chmod 755 plugins/*

	local cmd="docker build -t appscode/$IMG:$TAG-$K8S_VER-ac ."
	echo $cmd; $cmd
}

docker_push() {
	docker_up $IMG:$TAG-$K8S_VER-ac
}

docker_release() {
	local cmd="docker push appscode/$IMG:$TAG-$K8S_VER-ac"
	echo $cmd; $cmd
}

binary_repo $@
