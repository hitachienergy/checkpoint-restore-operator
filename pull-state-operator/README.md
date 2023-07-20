# pull-state-operator
The pull state operator implements pull state failover with [the example count application](../test-app). More specifically, it assumes that the application exposes an API allowing to retrieve and restore its state. The operator periodically pulls the state from the application and transfers it to the target node.
In case of a fault, the application is started on the target node and the operator calls the restore API to restore the state.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Application prerequisites
The pull-state-operator monitors appications having the following annotations:
```
    hitachienergy.com/ps_mode: "http-sf" # Ensures that the Liveness Handler deletes the old Pod and create a new Pod as soon as the fault is detected.
    hitachienergy.com/ps_interval: "10" # Pull interval
    hitachienergy.com/ps_path: "/state" # API endpoint exposed by the application to retrieve its state
    hitachienergy.com/ps_port: "8080" # Port number
```

### Running on the cluster

1. Set the following two environment variables:

```sh
export IMAGE_REGISTRY=<your-container-image-registry>
export CONTROL_PLANE_HOSTNAME=<k8s-ctrl-plane-node>
```

2. Perform the following replacements based on the environment variables set in the previous step:

```sh
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" helper/config/helper.daemonset.yaml
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" operator/config/manager/kustomization.yaml
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" operator/config/manager/manager.yaml
sed -i "s/CTRLHOSTNAME/${CONTROL_PLANE_HOSTNAME}/g" operator/config/manager/manager.yaml
```

3. Build and push your image:
	
```sh
cd operator
make docker-build docker-push
```
	
4. Deploy the controller to the cluster:

```sh
make deploy
```

5. Build the helper

```sh
cd helper
make docker-build docker-push
```

6. Deploy the helper

```sh
make deploy
```

### Viewing logs

- To view the operator logs, you can run: 

```
kubectl logs $(kubectl get pods --selector control-plane=controller-manager -n pull-state-operator-system -o jsonpath='{.items[0].metadata.name}') -n pull-state-operator-system
```

- To view the logs of the helper pods, replace `<your-k8s-node-hostname>` with the name of the Kubernetes node on which the helper is running and execute the following command:

```
kubectl logs $(kubectl get pods --selector app=ps-helper --field-selector spec.nodeName=<your-k8s-node-hostname> -n pull-state-operator-system -o jsonpath='{.items[0].metadata.name}') -n pull-state-operator-system
```

### Clean-up
To delete all resources related to the pull state operator from the cluster, you can run:

```sh
make undeploy
```