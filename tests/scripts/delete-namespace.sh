#!/bin/bash

# Verificar si se proporcionó el namespace
if [ "$#" -ne 1 ]; then
    echo "Uso: $0 <namespace>"
    exit 1
fi

NAMESPACE=$1


pods=$(kubectl  --kubeconfig $KUBE_CONFIG get pods -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}')

for pod in $pods; do
    echo "Borrando pod $pod en namespace $NAMESPACE..."
    kubectl  --kubeconfig $KUBE_CONFIG delete pod $pod -n $NAMESPACE --force --grace-period=0
done

echo "Todos los pods en $NAMESPACE han sido borrados."
