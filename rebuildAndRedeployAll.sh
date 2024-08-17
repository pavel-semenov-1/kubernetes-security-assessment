#!/bin/bash

docker build docker/demo-users/backend -t demo-users/backend &
docker build docker/demo-users/frontend -t demo-users/frontend &
docker build docker/demo-subjects/backend -t demo-subjects/backend &
docker build docker/demo-subjects/frontend -t demo-subjects/frontend &

wait

kubens demo-users && kubectl rollout restart deployment/demo-users-frontend deployment/demo-users-backend
kubens demo-subjects && kubectl rollout restart deployment/demo-subjects-frontend deployment/demo-subjects-backend