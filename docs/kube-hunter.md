---
title: Kube-hunter
author: Pavel Semenov
date: February 7, 2023
geometry: margin=1cm
---
# Kube-hunter

*Kube-hunter* is one of the most popular Kubernetes security scanners. This tool allows us to find vulnerabilities in cluster and its resources. There are three ways to run *kube-hunter*:

1. In a pod within a cluster. This indicates how exposed a cluster would be if one of the application pods is compromised (through a software vulnerability, for example).
2. Remotely from other machine using the IP address or domain name of the Kubernetes cluster. This will give you an attackers-eye-view of your Kubernetes setup.
3. Directly on a machine in the cluster.

## Preparation

To be able to do the experiments with *kube-hunter* we have to clone its official repository first.
``` {.bash}
git clone https://github.com/aquasecurity/kube-hunter.git
```
or
``` {.bash}
git clone git@github.com:aquasecurity/kube-hunter.git
```
Then `cd` into the cloned directory.

## Running kube-hunter in a pod

Let's now perform a security scan from within a pod on our prepared cluster.

NOTE: Hereafter we are using alias `k="minikube kubectl -- "` to simplify the cluster management.

First, switch to our test namespace (if not there already).
``` {.bash}
k config set-context --current --namespace=security-test
```

Then, we create a job, which will use *kube-hunter* to do the security scans.
``` {.bash}
k apply -f job.yaml
```

We can check state of the job by
``` {.bash}
k get po
```

When it is in state **Completed**, the results of the scan can be printed using the command below and the job's pod name
``` {.bash}
k logs <pod-name>
```

In my case the following report was generated:

``` {.txt}
Nodes
+-------------+------------+
| TYPE        | LOCATION   |
+-------------+------------+
| Node/Master | 10.244.0.1 |
+-------------+------------+
| Node/Master | 10.96.0.1  |
+-------------+------------+


Detected Services
+-------------+------------------+----------------------+
| SERVICE     | LOCATION         | DESCRIPTION          |
+-------------+------------------+----------------------+
| Kubelet API | 10.244.0.1:10250 | The Kubelet is the   |
|             |                  | main component in    |
|             |                  | every Node, all pod  |
|             |                  | operations goes      |
|             |                  | through the kubelet  |
+-------------+------------------+----------------------+
| API Server  | 10.96.0.1:443    | The API server is in |
|             |                  | charge of all        |
|             |                  | operations on the    |
|             |                  | cluster.             |
+-------------+------------------+----------------------+

Vulnerabilities
For further information about a vulnerability, search its ID in: 
https://avd.aquasec.com/
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
| ID     | LOCATION             | MITRE CATEGORY       | VULNERABILITY        | DESCRIPTION          | EVIDENCE             |
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
| None   | Local to Pod (kube-  | Lateral Movement //  | CAP_NET_RAW Enabled  | CAP_NET_RAW is       |                      |
|        | hunter-z2gtl)        | ARP poisoning and IP |                      | enabled by default   |                      |
|        |                      | spoofing             |                      | for pods.            |                      |
|        |                      |                      |                      |     If an attacker   |                      |
|        |                      |                      |                      | manages to           |                      |
|        |                      |                      |                      | compromise a pod,    |                      |
|        |                      |                      |                      |     they could       |                      |
|        |                      |                      |                      | potentially take     |                      |
|        |                      |                      |                      | advantage of this    |                      |
|        |                      |                      |                      | capability to        |                      |
|        |                      |                      |                      | perform network      |                      |
|        |                      |                      |                      |     attacks on other |                      |
|        |                      |                      |                      | pods running on the  |                      |
|        |                      |                      |                      | same node            |                      |
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
| KHV002 | 10.96.0.1:443        | Initial Access //    | K8s Version          | The kubernetes       | v1.25.3              |
|        |                      | Exposed sensitive    | Disclosure           | version could be     |                      |
|        |                      | interfaces           |                      | obtained from the    |                      |
|        |                      |                      |                      | /version endpoint    |                      |
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
| KHV005 | 10.96.0.1:443        | Discovery // Access  | Access to API using  | The API Server port  | b'{"kind":"APIVersio |
|        |                      | the K8S API Server   | service account      | is accessible.       | ns","versions":["v1" |
|        |                      |                      | token                |     Depending on     | ],"serverAddressByCl |
|        |                      |                      |                      | your RBAC settings   | ientCIDRs":[{"client |
|        |                      |                      |                      | this could expose    | CIDR":"0.0.0.0/0","s |
|        |                      |                      |                      | access to or control | ...                  |
|        |                      |                      |                      | of your cluster.     |                      |
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
| None   | Local to Pod (kube-  | Credential Access // | Access to pod's      | Accessing the pod's  | ['/var/run/secrets/k |
|        | hunter-z2gtl)        | Access container     | secrets              | secrets within a     | ubernetes.io/service |
|        |                      | service account      |                      | compromised pod      | account/namespace',  |
|        |                      |                      |                      | might disclose       | '/var/run/secrets/ku |
|        |                      |                      |                      | valuable data to a   | bernetes.io/servicea |
|        |                      |                      |                      | potential attacker   | ...                  |
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
| KHV050 | Local to Pod (kube-  | Credential Access // | Read access to pod's | Accessing the pod    | eyJhbGciOiJSUzI1NiIs |
|        | hunter-z2gtl)        | Access container     | service account      | service account      | ImtpZCI6IjBpR2tFeFFh |
|        |                      | service account      | token                | token gives an       | am1Sc01NWnU0TEJTMUhH |
|        |                      |                      |                      | attacker the option  | cHo1TEtLRjBNdEpHYXFk |
|        |                      |                      |                      | to use the server    | d0xXSTgifQ.eyJhdWQiO |
|        |                      |                      |                      | API                  | ...                  |
+--------+----------------------+----------------------+----------------------+----------------------+----------------------+
```

## Running kube-hunter remotely

TODO: finish this, need to find a way to access cluster remotely, because Minikube does not expose Kubernetes ports by default...

## Running kube-hunter on the cluster machine

To ssh into our "cluster machine" (which is basically a docker container) we can use the predefined minikube command
``` {.bash}
minikube ssh
```
First thing to do inside the cluster machine is to update the repository list.
``` {.bash}
sudo apt update -y
```
Then we can install *kube-hunter*. **Python3** should be already present on the machine, but there is no **pip** installed. 
``` {.bash}
sudo apt install python3-pip
```
Now we can use **pip** to install *kube-hunter*. Do not forget to add the installed binaries to the *PATH*.
``` {.bash}
python3 -m pip install kube-hunter
export PATH=$PATH:$HOME/.local/bin/
```
Running `kube-hunter` command will start the program. When prompted, choose to scan the ports of the local machine (since we are on the machine that runs the Kubernetes cluster). In my case the following report was produced:
``` {.txt}
Nodes
+-------------+--------------+
| TYPE        | LOCATION     |
+-------------+--------------+
| Node/Master | 192.168.49.2 |
+-------------+--------------+
| Node/Master | 10.244.0.1   |
+-------------+--------------+
| Node/Master | 10.244.0.1   |
+-------------+--------------+

Detected Services
+-------------+--------------------+----------------------+
| SERVICE     | LOCATION           | DESCRIPTION          |
+-------------+--------------------+----------------------+
| Kubelet API | 192.168.49.2:10250 | The Kubelet is the   |
|             |                    | main component in    |
|             |                    | every Node, all pod  |
|             |                    | operations goes      |
|             |                    | through the kubelet  |
+-------------+--------------------+----------------------+
| Kubelet API | 10.244.0.1:10250   | The Kubelet is the   |
|             |                    | main component in    |
|             |                    | every Node, all pod  |
|             |                    | operations goes      |
|             |                    | through the kubelet  |
+-------------+--------------------+----------------------+
| Kubelet API | 10.244.0.1:10250   | The Kubelet is the   |
|             |                    | main component in    |
|             |                    | every Node, all pod  |
|             |                    | operations goes      |
|             |                    | through the kubelet  |
+-------------+--------------------+----------------------+
| Etcd        | 192.168.49.2:2379  | Etcd is a DB that    |
|             |                    | stores cluster's     |
|             |                    | data, it contains    |
|             |                    | configuration and    |
|             |                    | current              |
|             |                    |     state            |
|             |                    | information, and     |
|             |                    | might contain        |
|             |                    | secrets              |
+-------------+--------------------+----------------------+

No vulnerabilities were found
```
Interestingly, no vulnerabilities were found in contrast to *in-pod* execution.

## Sources
- kube-hunter GitHub page [https://github.com/aquasecurity/kube-hunter](https://github.com/aquasecurity/kube-hunter)