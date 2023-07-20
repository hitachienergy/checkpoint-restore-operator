# cr-operator
This repository contains the work for the Master Thesis "[State Preserving Container Orchestration in Failover Scenarios](https://eldorado.tu-dortmund.de/handle/2003/41270)" written by Henri Schmidt (2023) in Cooperation with Hitachi Energy at the Research Center in Baden Dättwill. Supervisors for this work were Prof. Dr. Dr. Klaus Tycho-Förster, Dr. Raphael Eidenbenz and Dr. Zeineb Rejiba. Please keep in mind this work is scientific work and is not easily deployed and is not production ready. 
The thesis proposes a Kubernetes operator that enables users to experiment with the [new Kubernetes checkpoint feature](https://kubernetes.io/blog/2022/12/05/forensic-container-checkpointing-alpha/) to perform checkpoint/restore for a failover use case.
Please refer to the thesis for explanation of the architecture and concepts used by the cr-operator.

## Description
This operator regularly takes a checkpoint of monitored applications and transfers them to another node.
Should the health check of the application fail, it restores the most recent checkpoint on the other node.

## Contents
This repo contains everything needed to run the cr-operator. In addition, it contains:
- a `test-app`
 folder hosting an example application that you can use to test the operator.
 - a `pull-state-operator` folder containing an alternative failover operator that we use as a baseline. Instead of using checkpoint/restore, the pull-state-operator relies on the fact that an application provides an API for retrieving and restoring the state.
 - `0001-customization.patch`: This is a patch for CRI-O that includes a small bugfix and disables the cleanup of the CRIU logs. The patch applies to [this specific CRI-O commit](https://github.com/cri-o/cri-o/commits/642f60c471b6746652b1671637cbd17a07da5fcf) that we used during our experiments. It may be that the patch is no longer needed by the time you test this code. 
 
 You can refer to the README in the respective folders for a more detailed description.

## Prerequisites
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

Your cluster needs a special setup to work with checkpoint restore.

- First, your cluster needs to run on the latest version of [CRI-O](https://cri-o.io/) as currently only this container runtime is supported.
- Second, you need to install [CRIU](https://criu.org/Main_Page) on all your cluster nodes.
- Third, you need to enable Checkpoint/Restore (C/R) in all systems. The necessary steps are described in the following.

### Configuring Kubelet
Go to your Kubelet config and enable the Feature Gate `ContainerCheckpoint`. 

You can do so by adding `--feature-gates=ContainerCheckpoint=true` to `/var/lib/kubelet/kubeadm-flags.env`. Then, make sure to restart the Kubelet as follows:

```
sudo systemctl daemon-reload && sudo systemctl restart kubelet
```

### Configuring CRI-O

In your CRI-O config (`/etc/crio/crio.conf`), enable checkpoint restore:
```
[crio.runtime]
enable_criu_support = true
# mitigates a Bug in CRI-O, maybe this is fixed by the time you try this
drop_infra_ctr = false
```

### Configuring CRIU
In your CRIU config (`/etc/criu/runc.conf`), you need to set the following options:
```
tcp-close
skip-in-flight
# This is needed to mitigate another bug and is only needed with cgroups v2
manage-cgroups=ignore
```
**Note:** If this file does not exist, you can create it yourself.

### Prerequisites for the target application

Add the following annotations to your deployment to test out the controller:
```yaml
hitachienergy.com/cr_mode: "cr"
hitachienergy.com/cr_interval: "10" # interval in seconds for checkpointing
```

Your deployment needs to include an HTTP liveness probe for the operator to work correctly!

## Running on the cluster

1. Set the following two environment variables:

```sh
export IMAGE_REGISTRY=<your-container-image-registry>
export CONTROL_PLANE_HOSTNAME=<ctrl-plane-node>
```

2. replace IMAGEREGISTRY with the content of the environment variable set in the previous step.

```sh
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" agent/agent.daemonset.yaml
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" config/manager/kustomization.yaml
sed -i "s/IMAGEREGISTRY/${IMAGE_REGISTRY}/g" config/manager/manager.yaml
sed -i "s/CTRLHOSTNAME/${CONTROL_PLANE_HOSTNAME}/g" config/manager/manager.yaml
```

3. Create a secret for the controller with the kubelet certificates:

```sh
kubectl create namespace cr-operator-system
kubectl -n cr-operator-system create secret generic kubelet-client-certs --from-file=client.crt=/etc/kubernetes/pki/apiserver-kubelet-client.crt --from-file=client.key=/etc/kubernetes/pki/apiserver-kubelet-client.key
```

4. Build and push the cr-operator container image:
	
```sh
make docker-build docker-push
```
	
5. Deploy the controller to the cluster:

```sh
make deploy
```

6. Build and push the agent image:

```sh
cd agent
make docker-build docker-push
```

7. Deploy the agent to the cluster:

```sh
make deploy
```


8. Configure metrics

You can configure the cr-operator to store checkpoint metrics in a PostgreSQL database with the `DATABASE` environment variable in `config/manager/manager.yaml`.
For this to work, you need to add `display-stats` to your CRIU configuration. If the `DATABASE` environment variable is empty, metric collection is disabled.
## Debugging

Checking the logs of the CR-operator and/or the agents is useful for debugging issues. 

You can the check the CR-operator pod logs as follows:
```
kubectl logs $(kubectl get pods --selector control-plane=controllers-manager -n cr-operator-system -o jsonpath='{.items[0].metadata.name}') -n cr-operator-system
```

You can check the agent pod logs as follows while replacing `<your-k8s-node-name>` with the hostname of the Kubernetes node on which the agent is running:
```
kubectl logs $(kubectl get pods --selector app=cr-agent --field-selector spec.nodeName=<your-k8s-node-name> -n cr-operator-system -o jsonpath='{.items[0].metadata.name}') -n cr-operator-system
```
## Clean-up

To delete the test app, run the following in the `test-app` directory:
```
make undeploy
```

To delete everything related to the CR-operator, run:
```
make undeploy
```
