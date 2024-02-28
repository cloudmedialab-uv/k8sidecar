#!/bin/bash

#SET UP VARS

mkdir -p data/$EXPERIMENT_NAME


SECONDS=0
for archivo in "deploy/filters-java/"*.tmp 
do

# SET UP SCENARIO

kubectl apply --kubeconfig $KUBE_CONFIG -f $archivo

sleep 10

kubectl apply --kubeconfig $KUBE_CONFIG -f "deploy/functions/function.yml.tmp"

nname="$(basename "$archivo")"
name="${nname%%.yml.tmp}"

export TIMES_FILE=times$name.json

sleep 10

bash scripts/fullscale-replica.sh ffmpeg-fn-v2 2> /dev/null
echo $N_EXPERIMENTS
for i in $(seq 1 $N_EXPERIMENTS) 
do  
    echo "STARTING REPLICAS Experiment $i"
    bash scripts/send-request.sh "" "" "" 2> /dev/null
    sleep 2 
done

sleep 20

kubectl delete --kubeconfig $KUBE_CONFIG -f "deploy/functions/function.yml.tmp"

kubectl delete --kubeconfig $KUBE_CONFIG -f $archivo

sleep 10

# GET DATA 

#python3 scripts/getCoolTime.py times$name.json  data/coolstart/$EXPERIMENT_NAME/$name.txt

curl -s "$UPLOAD_SERVER_URL/$TIMES_FILE" > data/$EXPERIMENT_NAME/$name.json

done

echo "Tiempo total del experimento $SECONDS"
