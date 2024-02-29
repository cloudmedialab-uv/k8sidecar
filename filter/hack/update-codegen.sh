#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source "../../code-generator/kube_codegen.sh"

kube::codegen::gen_helpers \
    --input-pkg-root filter/pkg/apis \
    --output-base "../../" \
    --boilerplate "./boilerplate.go.txt"

kube::codegen::gen_client \
    --with-watch \
    --input-pkg-root filter/pkg/apis \
    --output-pkg-root filter/pkg/generated \
    --output-base "../../" \
    --boilerplate "./boilerplate.go.txt"
