export VMS="worker13-k8s-gpu-fat"
#NUMERO DE PETICIONES DE CADA PRUEBA
export N=4
#Number of Experiments
export N_EXPERIMENTS=100
# NOMBRE del fichero de tiempos
export TIMES_FILE=times.json
# Nombre del fichero de metricas gpu
export METRICS_FILE=gpu.csv
# Experiment with GPU
export GPU=true
# ACTIVATE METRICS SERVER FLAGS
export METRICS_FLAG=false

export MIN_REPLICAS=0
export MAX_REPLICAS=4

export CONCURRENCY=1

export DEVICES_PER_REPLICA=1

export FFMPEG_FLAGS="-c:v h264_nvenc -b:v 11M -maxrate 11M -minrate 11M -rc vbr -preset medium"