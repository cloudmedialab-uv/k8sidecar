#!/bin/bash

TAG="1.9.4.test"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker build $SCRIPT_DIR -t sidecar/filter/controller:$TAG -f $SCRIPT_DIR/deploy/docker/controller/Dockerfile

docker tag sidecar/filter/controller:$TAG routerdi1315.uv.es:33443/sidecar/filter/controller:$TAG


docker build $SCRIPT_DIR -t sidecar/filter/admission:$TAG -f $SCRIPT_DIR/deploy/docker/admission/Dockerfile

docker tag sidecar/filter/admission:$TAG routerdi1315.uv.es:33443/sidecar/filter/admission:$TAG