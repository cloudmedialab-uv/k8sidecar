# CONFIGURACIÓN KUBECTL MAQUINA MASTER
export KUBE_CONFIG=/home/tramuntana/.kube/config
# URL + PUERTO del servidor de subida de video
export UPLOAD_SERVER_URL=http://192.168.0.103:8082/upload
# URL + PUERTO del servidor de bajada de video
export DOWNLOAD_SERVER_URL=http://192.168.0.103:8081/VideosPeticiones
# URL del proxy de kubernetes
export K8S_URL=http://192.168.122.250:30080/video-coding