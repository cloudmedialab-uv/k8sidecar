#!/bin/bash
EXPERIMENTS="S1GPU S1NOGPU S2GPU S2NOGPU S3GPU S3NOGPU"
#EXPERIMENTS="S1GPU S1NOGPU"


source vars/vars.sh
bash inyection.sh create

find deploy/default -name "*.tmp" -exec kubectl apply --kubeconfig $KUBE_CONFIG -f {} \;

for EXPERIMENT_NAME in $EXPERIMENTS
do
    export EXPERIMENT_NAME=$EXPERIMENT_NAME
    echo "Ejecutando experimento: $EXPERIMENT_NAME"

    source vars/vars.sh
    bash inyection.sh create

    bash scripts/vm-manager.sh up $VMS

    bash tests/coldStart.sh

    bash inyection.sh clear
done

bash scripts/vm-manager.sh stop $VMS

echo "FINISH EXPERIMENTS"