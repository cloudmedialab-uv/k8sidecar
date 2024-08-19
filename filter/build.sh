#!/bin/bash

TAG="1.0.0"

REPO_CONTROLLER="cloudmedialab/filter-controller"
REPO_ADMISSION="cloudmedialab/filter-admission"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker build $SCRIPT_DIR -t sidecar/filter/controller:$TAG -f $SCRIPT_DIR/deploy/docker/controller/Dockerfile

docker tag sidecar/filter/controller:$TAG ${REPO_CONTROLLER}:${TAG}

docker push ${REPO_CONTROLLER}:${TAG}

docker build $SCRIPT_DIR -t sidecar/filter/admission:$TAG -f $SCRIPT_DIR/deploy/docker/admission/Dockerfile

docker tag sidecar/filter/admission:$TAG ${REPO_ADMISSION}:${TAG}

docker push ${REPO_ADMISSION}:${TAG}