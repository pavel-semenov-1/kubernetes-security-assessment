#!/bin/bash

kubectl -n ksa delete job kube-bench-runner prowler-runner trivy-runner
helm -n ksa uninstall ksa
sleep 2
helm -n ksa upgrade --install ksa .