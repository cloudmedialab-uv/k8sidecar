#!/bin/bash

EXPERIMENTS="latency"

source vars/vars.sh
bash inyection.sh create

find deploy/default -name "*.tmp" -exec kubectl apply --kubeconfig $KUBE_CONFIG -f {} \;

for EXPERIMENT_NAME in $EXPERIMENTS
do
    export EXPERIMENT_NAME=$EXPERIMENT_NAME
    echo "Ejecutando experimento: $EXPERIMENT_NAME"
    
    source vars/vars.sh
    bash inyection.sh create

    bash tests/latency.sh

    sleep 5

    bash inyection.sh clear
done

echo "FINISH EXPERIMENTS"
