#!/bin/bash

#SET UP VARS

export KUBE_CONFIG=/home/tramuntana/.kube/config-cluster

mkdir -p data/coolstart-envoy/envoy-8

SECONDS=0

#kubectl apply --kubeconfig $KUBE_CONFIG -f "deploy/filters-envoy/1.yml.tmp"

#sleep 15

kubectl apply --kubeconfig $KUBE_CONFIG -f "deploy/functions/function.yml.tmp"

name="1"

export TIMES_FILE=times$name.json

sleep 5

bash scripts/downscale-replica.sh ffmpeg-fn-v2 > /dev/null

for i in $(seq 1 $N_EXPERIMENTS) 
do  
    echo "STARTING REPLICAS Experiment $i"
    bash scripts/send-request.sh "" "" "" 2> /dev/null

    sleep 5

    bash scripts/downscale-replica.sh ffmpeg-fn-v2
done

kubectl delete --kubeconfig $KUBE_CONFIG -f "deploy/functions/function.yml.tmp" 2> /dev/null

#kubectl delete --kubeconfig $KUBE_CONFIG -f "deploy/filters-envoy/1.yml.tmp"

python3 scripts/getColdTime.py times$name.json  data/coolstart/$EXPERIMENT_NAME/$name.txt

echo "Tiempo total del experimento $SECONDS"
