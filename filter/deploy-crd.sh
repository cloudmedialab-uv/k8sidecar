#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

kubectl apply -f $SCRIPT_DIR/deploy/k8s/filter-controller.yaml
