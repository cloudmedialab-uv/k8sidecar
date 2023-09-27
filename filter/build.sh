
TAG="1.9.2.test"

docker build . -t sidecar/filter/controller:$TAG -f deploy/docker/controller/Dockerfile

docker tag sidecar/filter/controller:$TAG routerdi1315.uv.es:33443/sidecar/filter/controller:$TAG

docker push routerdi1315.uv.es:33443/sidecar/filter/controller:$TAG

docker build . -t sidecar/filter/admission:$TAG -f deploy/docker/admission/Dockerfile

docker tag sidecar/filter/admission:$TAG routerdi1315.uv.es:33443/sidecar/filter/admission:$TAG

docker push routerdi1315.uv.es:33443/sidecar/filter/admission:$TAG