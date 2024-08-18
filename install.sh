#!/bin/bash

docker build docker/demo-users/backend -t demo-users/backend
docker build docker/demo-users/frontend -t demo-users/frontend
docker build docker/demo-subjects/backend -t demo-subjects/backend
docker build docker/demo-subjects/frontend -t demo-subjects/frontend

kubectl create ns demo-users 
# 1.2. Pod Security Admission
kubectl label ns demo-users pod-security.kubernetes.io/enforce=privileged
kubectl create ns demo-subjects

kubectl -n demo-users create secret generic datasource-secret --from-literal=password=welcome1
kubectl -n demo-subjects create secret generic datasource-secret --from-literal=password=welcome1

helm -n demo-users install demo-users k8s/demo-users
helm -n demo-subjects install demo-subjects k8s/demo-subjects