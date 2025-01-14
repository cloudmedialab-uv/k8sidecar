#!/bin/bash

#SET UP VARS

export KUBE_CONFIG=/home/tramuntana/.kube/config-cluster

SECONDS=0

archivo=deploy/filters-envoy/10.yml.tmp

kubectl apply --kubeconfig $KUBE_CONFIG -f $archivo

sleep 10

kubectl apply --kubeconfig $KUBE_CONFIG -f "deploy/functions/function.yml.tmp"

name="10"

export TIMES_FILE=times$name.json

sleep 10

#bash scripts/fullscale-replica.sh ffmpeg-fn-v2 2> /dev/null

echo $N_EXPERIMENTS
for i in $(seq 1 $N_EXPERIMENTS) 
do  
    echo "Send Experiment $i"
    bash scripts/send-request.sh "" "" "" 
    sleep 2 
done

sleep 5

#kubectl delete --kubeconfig $KUBE_CONFIG -f "deploy/functions/function.yml.tmp"

#kubectl delete --kubeconfig $KUBE_CONFIG -f $archivo

#sleep 10

curl -s "$UPLOAD_SERVER_URL/$TIMES_FILE" > data/latency-envoy/10/times.json


echo "Tiempo total del experimento $SECONDS"
