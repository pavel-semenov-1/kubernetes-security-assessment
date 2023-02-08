---
title: Minikube
author: Pavel Semenov
date: February 7, 2023
geometry: margin=1cm
---
# Minikube

I am using probably the most common local Kubernetes implementation called Minikube.
It povides a few convenient tools for cluster and deployment management, including bundled Kubernetes Dashboard
and addns manager.
It is shipped with vast comprehensive documentation and has a huge global support community.
This installation is well-documented for every supported platform.

## Docker installation

Minikube requires container or Virtual Machine environment. I am using the Docker container environment.
When using Minikube with Docker there are two options: either use the default Docker daemon, which requires root priviliges or install rootless Docker. Since we are considering security issues in this thesis, I have chosen the latter to run Docker as current user. Installing rootless Docker comes down to the following three commands:
```
curl -o rootless-install.sh -fsSL https://get.docker.com/rootless
sh rootless-install.sh
export PATH=$HOME/bin:$PATH
```

However, Minikube also requires *cgroup v2* to be enabled. If the file `/sys/fs/cgroup/cgroup.controllers` is present on your machine, then you have it enabled. If not, refer to [this document](https://rootlesscontaine.rs/getting-started/common/cgroup2/).

We finish installation by executing these lines:
```
dockerd-rootless-setuptool.sh install -f
docker context use rootless
```

## Minikube installation

Since I am using Arch linux, and the `minikube` package is available in the community repository, I installed Minikube using the default package manager:
```
sudo pacman -S minikube
```

Then, assuming the rootless Docker is installed correctly, the cluster can be started via
```
minikube start --driver=docker --container-runtime=containerd
```

## Cluster setup

NOTE: Hereafter we are using alias `k="minikube kubectl -- "` to simplify the cluster management.

First, we have to create a separate namespace for our security tests.
```
k create ns security-test
```

Then, we change context to use this namespace
```
k config set-context --current --namespace=security-test
```

For convnience I have added the following function to my *.zshrc*:
```
kns() {
    kubectl config set-context --current --namespace=$1
}
```
This allows us to switch namespace simply by
```
kns security-test
```

## Provision deployment

We need some resources to perform the security tests. For now, I am using simple *nginx* server deployment.
```
k create deploy --image=nginx nginx
```
This will create the deployment with one replica. You can check if the pod (actual work unit) has started by 
```
k get po
```
Pod should be in state **Running**.

Then, we need to expose our pod through the service to be able to connect to it. List the pods again and copy the randomly generated name of the pod.
```
k expose <pod-name> --type=NodePort --port=80
```
You can then check the created service by executing
```
k get svc
```
Notice the assigned *NodePort*. We can now use it to access our deployment. First run `minikube ip` to get the IP address of our cluster. Then, use this IP and port number in your browser to access nginx server. If this doesn't work for you then use
```
minikube service <service-name> -n security-test
```

## Sources
- Minikube start documentation https://minikube.sigs.k8s.io/docs/start/
- Minikube Docker documentation https://minikube.sigs.k8s.io/docs/drivers/docker/
- Minikube documentation on accessing the applications https://minikube.sigs.k8s.io/docs/handbook/accessing/