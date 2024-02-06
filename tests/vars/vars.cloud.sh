# CONFIGURACIÓN KUBECTL MAQUINA MASTER
export KUBE_CONFIG=/home/tramuntana/.kube/config-cluster-k8s
# URL + PUERTO del servidor de subida de video
export UPLOAD_SERVER_USER=clouduser
export UPLOAD_SERVER_PRIVATE_KEY=/home/tramuntana/.keys/k8s-ssh-keys
export UPLOAD_SERVER_PATH=/home/clouduser/services/upload/data
export UPLOAD_SERVER_IP="192.168.0.242"
export UPLOAD_SERVER_URL=http://192.168.0.242:8080/upload

# URL + PUERTO del servidor de bajada de video
export DOWNLOAD_SERVER_URL=http://192.168.0.50:8080
export DOWNLOAD_SERVER_VIDEO_PATH=
# URL del proxy de kubernetes
export K8S_URL=http://192.168.0.116:30080/video-coding
# URL + PUERTO del servidor donde sube las metricas de gpu el sidecar
export GPU_METRICS_SERVER_URL=http://192.168.0.242:8080/upload
export METRICS_SERVER_URL=http://192.168.0.116:30090/stats