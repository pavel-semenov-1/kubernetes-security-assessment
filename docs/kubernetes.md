---
author: Pavel Semenov
date: 28.03.2023
title: Kubernetes
---

Kubernetes is a container orchestration platform. This chapter gives a quick overview of the Kubernetes and explains its advantages compared to the old-school virtual machine (hereafter VM) infrastructure.

## Architecture

Kubernetes requires at least three nodes to run. The master node is called **control plane**.  The **control plane** manages the **worker nodes** and the Pods in the cluster. In production environments, the control plane usually runs across multiple computers and a cluster usually runs multiple nodes, providing fault-tolerance and high availability. **Worker nodes** host the actual workload inside the cluster.

![Kubernetes architecture diagram](img/components-of-kubernetes.svg){width=100%}

### Control plane

**Control plane** runs the following components:
- kube-apiserver
- etcd
- kube-scheduler
- kube-controller-manager
- cloud-controller-manager

**Kube-apiserver** basically exposes the Kubernetes API, which is acting as a frontend for the Kubernetes control plane. **Etcd** is a key-value store, where all of the cluster data is stored. Each time a new pod is created, it is passed to the **kube-scheduler**, which assigns the pod to the specific node to run on (based on individual and collective resource requirements, hardware/software/policy constraints, affinity and anti-affinity specifications, data locality, inter-workload interference, and deadlines). Each Kubernetes resource has its own controller (e.g. *NodeController*, *JobController*, *ServiceAccountController*); all of them are compiled as one binary called **kube-controller-manager**. **Cloud-controller-manager** embeds cloud-specific control logic. It differs depending on the cloud provider or can be absent completely, when running Kubernetes locally.

### Worker nodes

Each **worker node** has a `kubelet` and `kube-proxy` installed. **Kubelet** is an agent that manages runnning
pods and containers. **Kube-proxy** is a network proxy that implements parts of the Kubernetes service concept.
It maintains network rules on nodes, making in- outside-cluster communitcation possible.

## Concepts & Resources

### Workloads
Minimal computing units in Kubernetes are *Containers*, which are running in *Pods*. However, to simplify the management of *Pods*, Kubernetes has workload resources, which manage the set of *Pods*. They make sure the desired number of *Pods* of right kind are running to match the declaration.

- *Deployments* and *ReplicaSets* are a good fit for stateless applications. Each pod in the *Deployment* is interchangebeable. *Deployments* are easily scalable and have built-in versioning and rollout mechanisms.
- *StatefulSet* allows to create sets of stateful applications. They might share the same *PersistentVolume* and replicate data between each other.
- *DaemonSet* defines Pods that provide node-local facilities. These might be fundamental to the operation of your cluster, such as a networking helper tool, or be part of an add-on.
- *Job* and *CronJob* define tasks that run to completion and then stop. *Jobs* represent one-off tasks, whereas *CronJobs* recur according to a schedule.

### Networking
Kubernetes networking model makes *Pods* look like VMs in the networking aspect. *Pods* on any nodes can communicate with each other without NAT. *Containers* inside the same *Pod* share its network meaning that they can reach each other using `localhost`.

Kubernetes networking addresses four concerns:

* *Containers* within a *Pod* use networking to communicate via loopback.
* Cluster networking provides communication between different *Pods*.
* The *Service* API lets you expose an application running in *Pods* to be reachable from outside your cluster.
    - *Ingress* provides extra functionality specifically for exposing HTTP applications, websites and APIs.
* You can also use *Services* to publish services only for consumption inside your cluster.

### Storage
Kubernetes does not ship a particular implementation of storage. However, it provides a range of resources that define the storage concept and supports different types of volumes. A Pod can use any number of volume types simultaneously. Ephemeral volume types have a lifetime of a pod, but persistent volumes exist beyond the lifetime of a pod. When a pod ceases to exist, Kubernetes destroys ephemeral volumes; however, Kubernetes does not destroy persistent volumes. For any kind of volume in a given pod, data is preserved across container restarts.

At its core, a volume is a directory, possibly with some data in it, which is accessible to the containers in a pod. How that directory comes to be, the medium that backs it, and the contents of it are determined by the particular volume type used.

*PersistentVolumes* and *PersistentVolumeClaims* are the resources kinds most important to understand here.

* A *PersistentVolume* (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. It is a resource in the cluster just like a node is a cluster resource. PVs are volume plugins like *Volumes*, but have a lifecycle independent of any individual Pod that uses the PV. This API object captures the details of the implementation of the storage, be that NFS, iSCSI, or a cloud-provider-specific storage system.

* A *PersistentVolumeClaim* (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources and PVCs consume PV resources. Pods can request specific levels of resources (CPU and Memory). Claims can request specific size and access modes (e.g., they can be mounted ReadWriteOnce, ReadOnlyMany or ReadWriteMany). *Once* and *Many* here refer to a number of *Pods* that can perform the read or write simultaneously.

## Comparison to VM
While Virtual Machines (VMs) provide a familiar way to manage infrastructure, they are slow and not flexible. That is, it might take as long as 60 seconds to boot up a VM. Booting a Kubernetes pod usually takes a few seconds. Migration is a lot more difficult on VMs. It is common to stumble upon a 20Gi VM image while the size of regular Docker image rarely goes above 5Gi. Even if the image is big, Docker implements multiple mechanisms for convenient image pushing and pulling making the migration fast and easy. 

![Virtual Machines vs Containers](img/containers-vs-virtual-machines.jpg){width=100%}

Among other advantages brought by cloud transformation are:

- Agile application creation and deployment: increased ease and efficiency of container image creation compared to VM image use.
- Continuous development, integration, and deployment: provides for reliable and frequent container image build and deployment with quick and efficient rollbacks (due to image immutability).
- Dev and Ops separation of concerns: create application container images at build/release time rather than deployment time, thereby decoupling applications from infrastructure.
- Observability: not only surfaces OS-level information and metrics, but also application health and other signals.
- Environmental consistency across development, testing, and production: runs the same on a laptop as it does in the cloud.
- Cloud and OS distribution portability: runs on Ubuntu, RHEL, CoreOS, on-premises, on major public clouds, and anywhere else.
- Application-centric management: raises the level of abstraction from running an OS on virtual hardware to running an application on an OS using logical resources.
- Loosely coupled, distributed, elastic, liberated micro-services: applications are broken into smaller, independent pieces and can be deployed and managed dynamically – not a monolithic stack running on one big single-purpose machine.
- Resource isolation: predictable application performance.
- Resource utilization: high efficiency and density.
- Service discovery and load balancing. Kubernetes can expose a container using the DNS name or using their own IP address. If traffic to a container is high, Kubernetes is able to load balance and distribute the network traffic so that the deployment is stable.
- Storage orchestration. Kubernetes allows you to automatically mount a storage system of your choice, such as local storages, public cloud providers, and more.
- Automated rollouts and rollbacks You can describe the desired state for your deployed containers using Kubernetes, and it can change the actual state to the desired state at a controlled rate. For example, you can automate Kubernetes to create new containers for your deployment, remove existing containers and adopt all their resources to the new container.
- Automatic bin packing. You provide Kubernetes with a cluster of nodes that it can use to run containerized tasks. You tell Kubernetes how much CPU and memory (RAM) each container needs. Kubernetes can fit containers onto your nodes to make the best use of your resources.
- Self-healing. Kubernetes restarts containers that fail, replaces containers, kills containers that don't respond to your user-defined health check, and doesn't advertise them to clients until they are ready to serve.
- Secret and configuration management Kubernetes lets you store and manage sensitive information, such as passwords, OAuth tokens, and SSH keys. You can deploy and update secrets and application configuration without rebuilding your container images, and without exposing secrets in your stack configuration.

There are, however, a few advantages of VMs compared to container infrastructure. One of the most important is security. Isolation of the VM from the host machine makes it less prone to attacks. We are talking about the security concepts of the cloud infrastructure and Kubernetes particularly in **Cloud secuity** chapter.

## Sources
- [https://kubernetes.io/docs/concepts/overview/components/](https://kubernetes.io/docs/concepts/overview/components/)
- [https://kubernetes.io/docs/concepts/workloads/](https://kubernetes.io/docs/concepts/workloads/)
- [https://kubernetes.io/docs/concepts/services-networking/](https://kubernetes.io/docs/concepts/services-networking/)
- [https://kubernetes.io/docs/concepts/storage/volumes/](https://kubernetes.io/docs/concepts/storage/volumes/)
- [https://kubernetes.io/docs/concepts/storage/persistent-volumes/](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
- [https://kubernetes.io/docs/concepts/security/overview/](https://kubernetes.io/docs/concepts/security/overview/)
- [https://www.weave.works/blog/a-practical-guide-to-choosing-between-docker-containers-and-vms](https://www.weave.works/blog/a-practical-guide-to-choosing-between-docker-containers-and-vms)
- [https://kubernetes.io/docs/concepts/overview/](https://kubernetes.io/docs/concepts/overview/)