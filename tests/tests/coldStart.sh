#!/bin/bash

#SET UP VARS

mkdir -p data/coolstart/$EXPERIMENT_NAME

FOLDER_SIDECAR=coolstartEncoder
FILE_FUNCTION=coolstartNOGPU


SECONDS=0
for archivo in "deploy/filters/$FOLDER_SIDECAR"/*.tmp 
do

# SET UP SCENARIO

kubectl apply --kubeconfig $KUBE_CONFIG -f $archivo

sleep 20

kubectl apply --kubeconfig $KUBE_CONFIG -f "deploy/functions/$FILE_FUNCTION.yml.tmp"

nname="$(basename "$archivo")"

if [[ $nname == *.yml.tmp ]]
then
    name="${nname%%.yml.tmp}"
else 
    name="${nname%%.yaml.tmp}"
fi

export TIMES_FILE=times$name.json

sleep 10

bash scripts/downscale-replica.sh ffmpeg-fn-v2 > /dev/null

for i in $(seq 1 $N_EXPERIMENTS) 
do  
    echo "STARTING REPLICAS Experiment $i"
    bash scripts/send-request.sh "-i prueba.mkv-c:v libx264 -an p.mp4" "\"$DOWNLOAD_SERVER_URL/prueba.mkv\"" "\"p.mp4\"" 2> /dev/null

    bash scripts/downscale-replica.sh ffmpeg-fn-v2
done

kubectl delete --kubeconfig $KUBE_CONFIG -f "deploy/functions/$FILE_FUNCTION.yml.tmp" 2> /dev/null

kubectl delete --kubeconfig $KUBE_CONFIG -f $archivo

# GET DATA 

python3 scripts/getCoolTime.py times$name.json  data/coolstart/$EXPERIMENT_NAME/$name.txt

done

echo "Tiempo total del experimento $SECONDS"
