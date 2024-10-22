---
title: Kube-hunter
author: Pavel Semenov
date: February 7, 2023
geometry: margin=1cm
---
# Trivy

*Trivy* is a comprehensive and versatile security scanner. *Trivy* has scanners that look for security issues, and targets where it can find those issues.

Targets (what Trivy can scan):

- Container Image
- Filesystem
- Git Repository (remote)
- Virtual Machine Image
- Kubernetes
- AWS

Scanners (what Trivy can find there):

- OS packages and software dependencies in use (SBOM)
- Known vulnerabilities (CVEs)
- IaC issues and misconfigurations
- Sensitive information and secrets
- Software licenses

## Installation

On my Mac machine I have installed Trivy using the Brew package manager:
``` {.bash}
brew install trivy
```

## Usage

To scan a Kubernetes cluster use the following command:

``` {.bash}
trivy k8s --report summary <context-name>
```

Here are the logs from my run:

``` {.txt}
2024-10-09T19:31:13+02:00	INFO	Node scanning is enabled
2024-10-09T19:31:13+02:00	INFO	If you want to disable Node scanning via an in-cluster Job, please try '--disable-node-collector' to disable the Node-Collector job.
237 / 237 [--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 1 p/s

Summary Report for rancher-desktop


Workload Assessment
┌───────────────┬───────────────────────────────────┬─────────────────────────┬────────────────────┬───────────────────┐
│   Namespace   │             Resource              │     Vulnerabilities     │ Misconfigurations  │      Secrets      │
│               │                                   ├───┬─────┬─────┬─────┬───┼───┬───┬───┬────┬───┼───┬───┬───┬───┬───┤
│               │                                   │ C │  H  │  M  │  L  │ U │ C │ H │ M │ L  │ U │ C │ H │ M │ L │ U │
├───────────────┼───────────────────────────────────┼───┼─────┼─────┼─────┼───┼───┼───┼───┼────┼───┼───┼───┼───┼───┼───┤
│ demo-users    │ StatefulSet/postgres              │ 8 │ 50  │ 63  │ 102 │   │   │ 1 │ 4 │ 9  │   │   │ 1 │   │   │   │
│ demo-users    │ Deployment/demo-users-frontend    │ 6 │ 118 │ 508 │ 489 │ 2 │   │ 2 │ 4 │ 10 │   │   │   │   │   │   │
│ demo-users    │ ConfigMap/datasource-config       │   │     │     │     │   │   │   │ 1 │    │   │   │   │   │   │   │
│ demo-users    │ Deployment/demo-users-backend     │   │ 44  │ 95  │ 17  │   │   │ 2 │ 8 │ 18 │   │   │   │   │   │   │
│ demo-subjects │ ConfigMap/datasource-config       │   │     │     │     │   │   │   │ 1 │    │   │   │   │   │   │   │
│ demo-subjects │ StatefulSet/postgres              │ 8 │ 50  │ 63  │ 102 │   │   │ 1 │ 4 │ 9  │   │   │ 1 │   │   │   │
│ demo-subjects │ Deployment/demo-subjects-frontend │ 6 │ 118 │ 508 │ 489 │ 2 │   │ 2 │ 4 │ 10 │   │   │   │   │   │   │
│ demo-subjects │ Deployment/demo-subjects-backend  │   │ 44  │ 95  │ 17  │   │   │ 2 │ 8 │ 18 │   │   │   │   │   │   │
└───────────────┴───────────────────────────────────┴───┴─────┴─────┴─────┴───┴───┴───┴───┴────┴───┴───┴───┴───┴───┴───┘
Severities: C=CRITICAL H=HIGH M=MEDIUM L=LOW U=UNKNOWN


Infra Assessment
┌─────────────┬──────────────────────────────────────────────┬────────────────────────┬─────────────────────┬───────────────────┐
│  Namespace  │                   Resource                   │    Vulnerabilities     │  Misconfigurations  │      Secrets      │
│             │                                              ├────┬─────┬─────┬───┬───┼───┬───┬────┬────┬───┼───┬───┬───┬───┬───┤
│             │                                              │ C  │  H  │  M  │ L │ U │ C │ H │ M  │ L  │ U │ C │ H │ M │ L │ U │
├─────────────┼──────────────────────────────────────────────┼────┼─────┼─────┼───┼───┼───┼───┼────┼────┼───┼───┼───┼───┼───┼───┤
│ kube-system │ DaemonSet/svclb-traefik-997e0b03             │    │     │ 40  │ 4 │   │   │ 7 │ 10 │ 18 │   │   │   │   │   │   │
│ kube-system │ DaemonSet/svclb-demo-backend-f5e6e704        │    │     │ 20  │ 2 │   │   │ 2 │ 6  │ 9  │   │   │   │   │   │   │
│ kube-system │ ConfigMap/extension-apiserver-authentication │    │     │     │   │   │   │   │ 1  │    │   │   │   │   │   │   │
│ kube-system │ Deployment/traefik                           │ 2  │  4  │ 42  │ 2 │   │   │   │ 2  │ 7  │   │   │   │   │   │   │
│ kube-system │ DaemonSet/svclb-demo-frontend-6aec7ab2       │    │     │ 20  │ 2 │   │   │ 2 │ 6  │ 9  │   │   │   │   │   │   │
│ kube-system │ Service/kube-dns                             │    │     │     │   │   │   │   │ 1  │    │   │   │   │   │   │   │
│ kube-system │ Deployment/local-path-provisioner            │ 1  │  2  │ 39  │ 2 │   │   │ 1 │ 4  │ 9  │   │   │   │   │   │   │
│ kube-system │ Job/helm-install-traefik                     │ 14 │ 104 │ 134 │ 4 │   │   │   │ 1  │ 6  │   │   │   │   │   │   │
│ kube-system │ Service/metrics-server                       │    │     │     │   │   │   │   │ 1  │    │   │   │   │   │   │   │
│ kube-system │ DaemonSet/svclb-demo-frontend-700a3f25       │    │     │ 20  │ 2 │   │   │ 2 │ 6  │ 9  │   │   │   │   │   │   │
│ kube-system │ DaemonSet/svclb-demo-backend-eebc0754        │    │     │ 20  │ 2 │   │   │ 2 │ 6  │ 9  │   │   │   │   │   │   │
│ kube-system │ Job/helm-install-traefik-crd                 │ 14 │ 104 │ 134 │ 4 │   │   │   │ 1  │ 6  │   │   │   │   │   │   │
│ kube-system │ Service/traefik                              │    │     │     │   │   │   │   │ 1  │    │   │   │   │   │   │   │
│ kube-system │ Deployment/coredns                           │ 3  │ 17  │ 22  │   │   │   │ 1 │ 4  │ 4  │   │   │   │   │   │   │
│ kube-system │ Deployment/metrics-server                    │    │     │     │   │   │   │   │ 2  │ 7  │   │   │   │   │   │   │
└─────────────┴──────────────────────────────────────────────┴────┴─────┴─────┴───┴───┴───┴───┴────┴────┴───┴───┴───┴───┴───┴───┘
Severities: C=CRITICAL H=HIGH M=MEDIUM L=LOW U=UNKNOWN


RBAC Assessment
┌─────────────┬────────────────────────────────────────────────────────────────────┬───────────────────┐
│  Namespace  │                              Resource                              │  RBAC Assessment  │
│             │                                                                    ├───┬───┬───┬───┬───┤
│             │                                                                    │ C │ H │ M │ L │ U │
├─────────────┼────────────────────────────────────────────────────────────────────┼───┼───┼───┼───┼───┤
│ kube-system │ Role/system:controller:bootstrap-signer                            │   │   │ 1 │   │   │
│ kube-system │ Role/system::leader-locking-kube-controller-manager                │   │   │ 1 │   │   │
│ kube-system │ Role/system::leader-locking-kube-scheduler                         │   │   │ 1 │   │   │
│ kube-system │ Role/system:controller:token-cleaner                               │   │   │ 1 │   │   │
│ kube-system │ Role/system:controller:cloud-provider                              │   │   │ 1 │   │   │
│ kube-public │ Role/system:controller:bootstrap-signer                            │   │   │ 1 │   │   │
│ demo-users  │ RoleBinding/demo-users-rolebinding                                 │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:endpointslicemirroring-controller    │   │ 1 │   │   │   │
│             │ ClusterRoleBinding/cluster-admin                                   │   │   │ 1 │   │   │
│             │ ClusterRole/k3s-cloud-controller-manager                           │ 1 │ 1 │ 1 │   │   │
│             │ ClusterRole/edit                                                   │ 2 │ 4 │ 6 │   │   │
│             │ ClusterRole/admin                                                  │ 3 │ 4 │ 6 │   │   │
│             │ ClusterRole/system:controller:root-ca-cert-publisher               │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:legacy-service-account-token-cleaner │ 1 │   │   │   │   │
│             │ ClusterRole/system:kube-scheduler                                  │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:horizontal-pod-autoscaler            │ 2 │   │   │   │   │
│             │ ClusterRole/system:controller:persistent-volume-binder             │ 1 │ 2 │ 1 │   │   │
│             │ ClusterRole/system:controller:resourcequota-controller             │ 1 │   │   │   │   │
│             │ ClusterRole/system:controller:statefulset-controller               │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:replication-controller               │   │   │ 2 │   │   │
│             │ ClusterRole/system:kube-controller-manager                         │ 5 │   │   │   │   │
│             │ ClusterRole/cluster-admin                                          │ 2 │   │   │   │   │
│             │ ClusterRole/system:node                                            │ 1 │   │ 1 │   │   │
│             │ ClusterRole/system:controller:replicaset-controller                │   │   │ 2 │   │   │
│             │ ClusterRoleBinding/helm-kube-system-traefik-crd                    │   │   │ 1 │   │   │
│             │ ClusterRole/system:aggregate-to-edit                               │ 2 │ 4 │ 6 │   │   │
│             │ ClusterRole/system:controller:endpointslice-controller             │   │ 1 │   │   │   │
│             │ ClusterRole/system:aggregate-to-admin                              │ 1 │   │   │   │   │
│             │ ClusterRole/traefik-kube-system                                    │ 1 │   │   │   │   │
│             │ ClusterRole/local-path-provisioner-role                            │ 1 │ 1 │ 1 │   │   │
│             │ ClusterRole/system:controller:endpoint-controller                  │   │ 1 │   │   │   │
│             │ ClusterRole/system:controller:pod-garbage-collector                │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:expand-controller                    │ 1 │   │   │   │   │
│             │ ClusterRole/system:controller:job-controller                       │   │   │ 2 │   │   │
│             │ ClusterRole/system:controller:generic-garbage-collector            │ 1 │   │   │   │   │
│             │ ClusterRole/system:controller:ttl-after-finished-controller        │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:deployment-controller                │   │   │ 3 │   │   │
│             │ ClusterRole/system:controller:node-controller                      │   │   │ 1 │   │   │
│             │ ClusterRoleBinding/helm-kube-system-traefik                        │   │   │ 1 │   │   │
│             │ ClusterRole/system:controller:cronjob-controller                   │   │   │ 3 │   │   │
│             │ ClusterRole/system:controller:namespace-controller                 │ 1 │   │   │   │   │
│             │ ClusterRole/system:controller:daemon-set-controller                │   │   │ 1 │   │   │
└─────────────┴────────────────────────────────────────────────────────────────────┴───┴───┴───┴───┴───┘
Severities: C=CRITICAL H=HIGH M=MEDIUM L=LOW U=UNKNOWN

```