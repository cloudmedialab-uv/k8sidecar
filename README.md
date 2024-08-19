# K8Sidecar

This repo contains the source code of the Filter Custom Resource Definition (CRD) that enables the
declaration and inyection of proxy sidecars.

Besides, it contains an example of the definition of two filters:
 - logauth filter with two proxy sidecars: one that performs logging and another that performs authentication
 - ratelimiter filter with one proxy sidecar
The example shows how these two filters can be used in a deployment to inyect the three proxy sidecars.

## Prerequisites

-   [Minikube](https://minikube.sigs.k8s.io/docs/start/) or your custom cluster
-   [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
-   [Docker](https://www.docker.com/products/docker-desktop)

## Installation

### 1. Install Minikube

To install Minikube, follow the instructions on the [official Minikube page](https://minikube.sigs.k8s.io/docs/start/).

Although we will use Minikube for testing and local development, our application is designed to be deployed on any Kubernetes cluster. Therefore, if you already have another Kubernetes environment set up and prefer to use it, feel free to do so, ensuring you follow the specific configuration instructions for your cluster.

### 2. Start Minikube

```bash
minikube start
```

### 3. Clone the Repository

```bash
git clone https://github.com/cloudmedialab-uv/k8sidecar.git
cd k8sidecar
```

You have two options:
 - Build the sources to create the two images that manage the new filter CRD, or
 - Use the images that are in dockerhub.

#### 3.1 Build the sources

If you choose to compile the images manually, edit the `build.sh` script to adapt the registry where the two images
will be stored and the tag of the images. Then run the script:

```bash
bash filter/build.sh
```

This script will create two images:
 - The container image of the controller for the filter
 - The container image of the webhook for the mutating admission controller
Finally, these two images will be uploaded to the registry.

These two Docker images will be built for later use.

#### 3.2 Use the already created images in dockerhub

If you prefer to use the available precompiled images, simply proceed with the deployment of the CRD, as these images are hosted in the official repository and will be downloaded automatically during deployment.


### Deploy the Filter Custom Resource Definition (CRD)

To deploy this CRD just run:

```bash
bash deploy-crd.sh
```

This just will deploy to kubernetes the `filter-controller.yaml` file located in folder `deploy/k8s`.

## Examples

In this section, we will use proxy sidecar images created with the Java and Go client libraries. You can find references to both libraries and the source code of the proxy sidecars in their respective GitHub repositories ([Java Client Library](https://github.com/cloudmedialab-uv/k8sidecar-java-lib), [Go Client Library](https://github.com/cloudmedialab-uv/k8sidecar-go-lib)).

### ratelimiter filter

In the first example, we will use a filter created with the Go library that acts as a rate limiter. You can find the reference to this example [here](https://github.com/cloudmedialab-uv/k8sidecar-go-lib/tree/main/examples/ratelimiter).

Below is the YAML file specification of this filter. This file is located in the `examples` folder of the repository.

```yaml
apiVersion: k8sidecar.io/v1
kind: Filter
metadata:
  name: ratelimiter
spec:
  sidecars:
    - image: "routerdi1315.uv.es:33443/sidecar/ratelimiter:1.0.0"
      name: "ratelimiter-container"
      priority: 2
```

To create this filter in kubernetes, run the following command:

```bash
kubectl apply -f examples/ratelimiter.yaml
```

To check that the filter has been created:
```bash
kubectl get filter
```

### logauth filter

The second example consists of proxy sidecars written with the Java library, which act as authentication and logging requests. References to these sidecars can be found [here for authentication](https://github.com/your-username/auuth-example) and [here for logging](https://github.com/cloudmedialab-uv/k8sidecar-java-lib/tree/main/examples/logging).

Below is the YAML specification of the filter that uses these two proxy sidecars:

```yaml
apiVersion: k8sidecar.io/v1
kind: Filter
metadata:
  name: logauth
spec:
  sidecars:
    - image: "routerdi1315.uv.es:33443/sidecar/logging:1.0.0"
      name: "logging-container"
      priority: 1
    - image: "routerdi1315.uv.es:33443/sidecar/auth:1.0.0"
      name: "auth-container"
      priority: 3
      env:
        - name: AUTH_TOKEN_NAME
          value: AUTH_TOKEN
        - name: AUTH_TOKEN_KEY
          value: "password"
```

To create this filter, run the following command:

```bash
kubectl apply -f examples/logauth.yaml
```

To check that the filter has been created:
```bash
kubectl get filter
```

## Run a deployment that uses these two filter

Once the filters are created, they can be injected into your deployments using Kubernetes labels. Reference the name of the desired filter to inject it into your deployment.

Below is an example of a deployment manifest. In this example, the `logauth: "sidecar"` and `ratelimiter: "sidecar"` labels are used to inject the proxy sidecars defined in these filters (the logging and authentication, and rate limiter sidecars, respectively).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-deployment
  annotations:
    k8sidecar.port: "PORT"
  labels:
    logauth: "sidecar" # request to inject logging and auth sidecars
    ratelimiter: "sidecar" # request to inject ratelimiter sidecar
spec:
  replicas: 1
  selector:
    matchLabels:
      app: echo
  template:
    metadata:
      labels:
        app: echo
    spec:
      containers:
      - name: echo-container
        image: ealen/echo-server:latest
        env:
        - name: PORT
          value: "80"
---
apiVersion: v1
kind: Service
metadata:
  name: echo-service
spec:
  selector:
    app: echo
  ports:
  - port: 80
    targetPort: 80
```

By adding these labels to the metadata of your deployment, the corresponding sidecars will be automatically injected into the pods created by this deployment, enhancing the functionality of the deployed application with logging, authentication, and rate limiting features (not present in the application image).

To create this deployment, run the following command:

```bash
kubectl create namespace filter-usage
kubectl apply -f examples/deployment.yaml
```

If everything is configured correctly, you should see the successful deployment messages in your terminal. To verify that the deployment and sidecar injection were successful, you can check the pods' statuses with the following command:

```bash
kubectl get pod -n filter-usage
```

A successful deployment will show the pod in a Running status with 4 containers, confirming that everything is operating as expected.

### Expected output

```
NAME                               READY   STATUS         RESTARTS   AGE
echo-deployment-788dffbb7f-mvqlm   4/4     Running        0          5s
```

### Testing
We will deploy a pod with `curl` to make HTTP requests:
```bash
kubectl run curl --image=curlimages/curl:8.5.0 --restart=Never --namespace filter-usage --command -- sleep 3600
```

We can check that the authentication proxy filters out the requests that do not provide the HTTP header configured in this sidecar:
```bash
kubectl exec curl -n filter-usage -- curl http://echo-service
```

Check that the authentication proxy passes the request if we provide the HTTP header configured in this sidecar:
```bash
kubectl exec curl -n filter-usage -- curl -H "AUTH_TOKEN: password" http://echo-service
```

### Clean

```bash
kubectl delete namespace filter-usage
```

If you do not need the filters you can delete them:
```bash
kubectl delete logauth
kubectl delete ratelimiter
```