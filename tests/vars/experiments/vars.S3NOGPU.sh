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
export GPU=false
# ACTIVATE METRICS SERVER FLAGS
export METRICS_FLAG=false

export MIN_REPLICAS=0
export MAX_REPLICAS=1

export CONCURRENCY=4

export FFMPEG_FLAGS="-c:v libx264 -b:v 11M -maxrate 11M -minrate 11M -preset medium"