#!/bin/bash

kubectl -n ksa scale deployment.apps/ksa-aggregator deployment.apps/ksa-dashboard deployment.apps/ksa-parser statefulset.apps/postgres --replicas=0
# wait grace period
sleep 30
helm -n ksa upgrade --install ksa .
kubectl -n ksa scale deployment.apps/ksa-aggregator deployment.apps/ksa-dashboard deployment.apps/ksa-parser statefulset.apps/postgres --replicas=1