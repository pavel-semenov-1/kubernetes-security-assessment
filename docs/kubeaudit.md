---
title: Kubeaudit
author: Pavel Semenov
date: February 7, 2023
geometry: margin=1cm
---
# Kubeaudit

**Kubeaudit** is a tool to audit Kubernetes clusters for various different security concerns, such as:
- run as non-root
- use a read-only root filesystem
- drop scary capabilities, don't add new ones
- don't run privileged

It has three modes:
1. Manifest mode
2. Local mode
3. Cluster mode
We will try out each mode and document our findings below.

## Preparation

## Manifest mode
**Kubeaudit** can be executed against a manifest file to scan it for security issues. The command is
```
kubeaudit all -f "/path/to/manifest.yml"
```
Running **kubeaudit** in manifest mode against our test setup produced the following output:
``` {.bash}
---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx

--------------------------------------------

-- [error] AppArmorAnnotationMissing
   Message: AppArmor annotation missing. The annotation 'container.apparmor.security.beta.kubernetes.io/nginx' should be added.
   Metadata:
      Container: nginx
      MissingAnnotation: container.apparmor.security.beta.kubernetes.io/nginx

-- [error] AutomountServiceAccountTokenTrueAndDefaultSA
   Message: Default service account with token mounted. automountServiceAccountToken should be set to 'false' on either the ServiceAccount or on the PodSpec or a non-default service account should be used.

-- [error] CapabilityOrSecurityContextMissing
   Message: Security Context not set. The Security Context should be specified and all Capabilities should be dropped by setting the Drop list to ALL.
   Metadata:
      Container: nginx

-- [warning] ImageTagMissing
   Message: Image tag is missing.
   Metadata:
      Container: nginx

-- [warning] LimitsNotSet
   Message: Resource limits not set.
   Metadata:
      Container: nginx

-- [error] RunAsNonRootPSCNilCSCNil
   Message: runAsNonRoot should be set to true or runAsUser should be set to a value > 0 either in the container SecurityContext or PodSecurityContext.
   Metadata:
      Container: nginx

-- [error] AllowPrivilegeEscalationNil
   Message: allowPrivilegeEscalation not set which allows privilege escalation. It should be set to 'false'.
   Metadata:
      Container: nginx

-- [warning] PrivilegedNil
   Message: privileged is not set in container SecurityContext. Privileged defaults to 'false' but it should be explicitly set to 'false'.
   Metadata:
      Container: nginx

-- [error] ReadOnlyRootFilesystemNil
   Message: readOnlyRootFilesystem is not set in container SecurityContext. It should be set to 'true'.
   Metadata:
      Container: nginx

-- [error] SeccompProfileMissing
   Message: Pod Seccomp profile is missing. Seccomp profile should be added to the pod SecurityContext.
```
As you can see a lot of misconfigurations were reported. Now let's ask **kubeaudit** to fix the issues for us.
We can do so buy running:
``` {.bash}
kubeaudit autofix -f "/path/to/manifest.yml" -o "/path/to/fixed"
```
The fixed manifest looks as follows:
``` {.bash}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  template:
    spec:
      containers:
        - image: nginx
          name: nginx
          resources: {}
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL
            privileged: false
            readOnlyRootFilesystem: true
            runAsNonRoot: true
      automountServiceAccountToken: false
      securityContext:
        seccompProfile:
          type: RuntimeDefault
    metadata:
      annotations:
        container.apparmor.security.beta.kubernetes.io/nginx: runtime/default
  selector: null
  strategy: {}
```
Now only a few warnings are left:
```
---------------- Results for ---------------

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx

--------------------------------------------

-- [warning] ImageTagMissing
   Message: Image tag is missing.
   Metadata:
      Container: nginx

-- [warning] LimitsNotSet
   Message: Resource limits not set.
   Metadata:
      Container: nginx
```

## Local mode

## Cluster mode

## Conclusion