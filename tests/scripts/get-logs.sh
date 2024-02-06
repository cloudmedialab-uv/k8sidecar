#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Uso: $0 <nombre-del-deployment>"
    exit 1
fi

DEPLOYMENT_NAME=$1

# Obtener los pods del despliegue
PODS=$(kubectl --kubeconfig $KUBE_CONFIG get pods -l=app=$DEPLOYMENT_NAME -o=jsonpath='{.items[*].metadata.name}')

for pod in $PODS; do
    echo "Obteniendo registros del pod $pod"
    # Obtener contenedores del pod
    CONTAINERS=$(kubectl --kubeconfig $KUBE_CONFIG get pod $pod -o=jsonpath='{.spec.containers[*].name}')
    for container in $CONTAINERS; do
        echo "Obteniendo registros del contenedor $container del pod $pod"
        kubectl --kubeconfig $KUBE_CONFIG logs $pod -c $container
    done
done
