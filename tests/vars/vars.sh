# NOMBRE DEL USUARIO DE LA MAQUINA CON LOS WORKERS
export REMOTE_USER=clouduser
# IP DE LA MAQUINA CON LOS WORKERS
export REMOTE_HOST="192.168.0.112"
# CERTIFICADO DE LA MAQUINA CON LOS WORKERS
export REMOTE_PRIVATE_KEY=/home/tramuntana/.keys/k8s-ssh-keys
# url del broker de entrada de peticiones
export BROKER_URL=http://broker-ingress.knative-eventing.svc.cluster.local/default/video-coding-broker
# Tiempo sin peticiones antes de que escale a 0 la replica
export SCALE_TO_ZERO_GRACE_PERIOD="15s"

# namespace donde se encuentran los recursos de kubernetes
export NAMESPACE=default

export DEPLOY_ENV="cloud"

if [ -f "vars/vars.${DEPLOY_ENV}.sh" ]; then
    source "vars/vars.${DEPLOY_ENV}.sh"
fi

if [ -f "vars/experiments/vars.${EXPERIMENT_NAME}.sh" ]; then
    source "vars/experiments/vars.${EXPERIMENT_NAME}.sh"
fi

