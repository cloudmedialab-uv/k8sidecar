kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.9.2/serving-crds.yaml

kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.9.2/serving-core.yaml


kubectl apply -f https://github.com/knative/net-contour/releases/download/knative-v1.9.3/contour.yaml

kubectl apply -f https://github.com/knative/net-contour/releases/download/knative-v1.9.3/net-contour.yaml

kubectl patch configmap/config-network \
 --namespace knative-serving \
 --type merge \
 --patch '{"data":{"ingress-class":"contour.ingress.networking.knative.dev"}}'


kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.9.7/eventing-crds.yaml

kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.9.7/eventing-core.yaml

kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.9.7/in-memory-channel.yaml

kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.9.7/mt-channel-broker.yaml

kubectl create configmap config-defaults --namespace=knative-serving --from-literal=revision-timeout-seconds="3600" --from-literal=max-revision-timeout-seconds="7200"


kubectl create broker video-coding-broker

kubectl create secret docker-registry routerdi-registry-creds \
    --docker-server routerdi1315.uv.es:33443 \
    --docker-username=cloudlab \
    --docker-password=registro

kubectl patch serviceaccount default  -p "{\"imagePullSecrets\": [{\"name\": \"routerdi-registry-creds\"}]}"
