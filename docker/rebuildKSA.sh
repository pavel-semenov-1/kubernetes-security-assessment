#!/bin/bash

docker build docker/ksa-aggregator -t ksa/aggregator
docker build docker/ksa-parser -t ksa/parser
docker build docker/ksa-dashboard -t ksa/dashboard
docker build docker/ksa-postgres -t ksa/postgres
docker build docker/prowler -t ksa/prowler
docker build docker/kubescape -t ksa/kubescape