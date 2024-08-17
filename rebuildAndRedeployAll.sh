#!/bin/bash

docker build docker/demo-users/backend -t demo-users/backend
docker build docker/demo-users/frontend -t demo-users/frontend
docker build docker/demo-subjects/backend -t demo-subjects/backend
docker build docker/demo-subjects/frontend -t demo-subjects/frontend

wait

kubectl -n demo-users rollout restart deployment/demo-users-frontend deployment/demo-users-backend
kubectl -n demo-subjects rollout restart deployment/demo-subjects-frontend deployment/demo-subjects-backend

sleep 60

docker container prune -f
docker image prune -af