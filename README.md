# K8Sidecar

Brief description of your application.

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
git clone https://github.com/tu-usuario/tu-repositorio.git
cd tu-repositorio
```

### 4. Build the Sources

To build the sources, you have two options: you can do it manually or use the images already available in the official repository of the application.

If you choose to compile the images manually, run the `build.sh` script provided in the repository:

```bash
bash filter/build.sh
```

This script will compile the necessary Docker images for the application and properly tag them for later use.

If you prefer to use the available precompiled images, simply proceed with the deployment of the application, as these images are hosted in the official repository and will be downloaded automatically during deployment.

### 5. Deploy the Application in Kubernetes

To deploy the application on your Kubernetes cluster, you can use the `start.sh ` script included in the repository:

```bash
bash filter/start.sh
```

This script automates the deployment process, applying the necessary Kubernetes configurations to install and run the application on your cluster. Ensure you have the appropriate credentials configured and are authenticated in your cluster before running the script. In case you are using Minikube, ensure that the cluster is running with `minikube start`.

## Run Examples

In this section, we will use example images created with the Java and Go client libraries. You can find references to both libraries in their respective GitHub repositories ([Java Client Library](https://github.com/your-username/java-client-library), [Go Client Library](https://github.com/your-username/go-client-library)).

### Example 1

In the first example, we will use a filter created with the Go library that acts as a rate limiter. You can find the reference to this example [here](https://github.com/your-username/ratelimiter-example).

Below is the YAML file description for the Go filter which will be installed in Kubernetes. This file is located in the `examples` folder of the repository.

```yaml
apiVersion: filtercontroller.ks.io/v1
kind: Filter
metadata:
    name: ratelimiter
spec:
    sidecars:
        - image: "routerdi1315.uv.es:33443/sidecar/ratelimiter:1.0.0"
          name: "ratelimiter_container"
```

To execute this example, run the following command:

```bash
kubectl apply -f examples/ratelimiter.yaml
```

### Example 2

The second example consists of two images written with the Java library, which act as authentication and logging request filters. References to these codes can be found [here for authentication](https://github.com/your-username/auuth-example) and [here for logging](https://github.com/your-username/log-example).

Below is the YAML file description for the Java filters, which will also be installed in Kubernetes:

```yaml
apiVersion: filtercontroller.ks.io/v1
kind: Filter
metadata:
    name: logauth
spec:
    sidecars:
        - image: "routerdi1315.uv.es:33443/sidecar/logging:1.0.0"
          name: "logging_container"
        - image: "routerdi1315.uv.es:33443/sidecar/auth:1.0.0"
          name: "auth_container"
          env:
              - name: AUTH_TOKEN_NAME
                value: AUTH_TOKEN
              - name: AUTH_TOKEN_KEY
                value: "password"
```

To execute this example, run the following command:

```bash
kubectl apply -f examples/logauth.yaml
```

## Run a Deployment

Once the filters are installed, they can be injected into your deployments using Kubernetes labels. Reference the name of the desired filter to inject it into your deployment.

Below is an example of a deployment manifest for an echo server. In this example, the `logauth: "sidecar"` and `ratelimiter: "sidecar"` labels are used to inject the logging and authentication, and rate limiter sidecars, respectively.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: echo-deployment
    labels:
        logauth: "sidecar" # inject logging and auth sidecars
        ratelimiter: "sidecar" # inject ratelimiter sidecar
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
                  ports:
                      - containerPort: 80
```

By adding these labels to the metadata of your deployment, the corresponding sidecars will be automatically injected into the pods created by this deployment, enhancing the functionality of the deployed application with logging, authentication, and rate limiting features.

Then apply the deployment, run the following command:

```bash
kubectl apply -f examples/deployment.yaml
```

If everything is configured correctly, you should see the successful deployment messages in your terminal. To verify that the deployment and sidecar injection were successful, you can check the pods' statuses with the following command:

```bash
kubectl get pods
```

A successful deployment will show the pod in a Running status with 4 containers, confirming that everything is operating as expected.

### Expected output

```
NAME                               READY   STATUS         RESTARTS   AGE
echo-deployment-788dffbb7f-mvqlm   4/4     Running        0          5s
```
