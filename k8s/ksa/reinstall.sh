#!/bin/bash

kubectl -n ksa delete job kube-bench-runner prowler-runner trivy-runner kubescape-runner
helm -n ksa uninstall ksa
sleep 5
helm -n ksa upgrade --install ksa .
