#!/bin/bash

NAMESPACE=$NAMESPACE

DEPLOYMENT_NAME="$1"

if [ -z "$DEPLOYMENT_NAME" ]; then
    echo "Por favor, proporciona el nombre del deployment."
    exit 1
fi

while true
do
    # Obtener el número de pods en estado "Running" asociados con el deployment
    RUNNING_PODS=$(kubectl --kubeconfig $KUBE_CONFIG get pods -n "$NAMESPACE" -l app="$DEPLOYMENT_NAME" -o=jsonpath='{.items[?(@.status.phase=="Running")].metadata.name}' | wc -w)

    if [ "$RUNNING_PODS" == "$MAX_REPLICAS" ]; then
        echo "El deployment $DEPLOYMENT_NAME ha escalado a $MAX_REPLICAS y todos los pods están en estado 'Running'."
        break
    fi
    
    sleep 5
done

