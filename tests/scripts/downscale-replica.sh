#!/bin/bash

NAMESPACE=$NAMESPACE

DEPLOYMENT_NAME="$1"

if [ -z "$DEPLOYMENT_NAME" ]; then
    echo "Por favor, proporciona el nombre del deployment."
    exit 1
fi

while true
do

        TERMINATING_PODS=$(kubectl --kubeconfig $KUBE_CONFIG get pods -n "$NAMESPACE" -l=app="$DEPLOYMENT_NAME" --no-headers | wc -l)

        if [ "$TERMINATING_PODS" == "0" ]; then
            echo "El deployment $DEPLOYMENT_NAME ha escalado a 0 y todos los pods han terminado."
            break
        fi
    sleep 5
done