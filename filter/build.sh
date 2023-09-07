docker build . -t sidecar/filter/controller:1.7.test -f deploy/docker/controller/Dockerfile

docker tag sidecar/filter/controller:1.7.test routerdi1315.uv.es:33443/sidecar/filter/controller:1.7.test

docker push routerdi1315.uv.es:33443/sidecar/filter/controller:1.7.test

docker build . -t sidecar/filter/admission:1.7.test -f deploy/docker/admission/Dockerfile

docker tag sidecar/filter/admission:1.7.test routerdi1315.uv.es:33443/sidecar/filter/admission:1.7.test

docker push routerdi1315.uv.es:33443/sidecar/filter/admission:1.7.test