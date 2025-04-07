#!/bin/bash

kubectl -n ksa delete job kube-bench-runner prowler-runner trivy-runner
helm -n ksa uninstall ksa
helm -n ksa upgrade --install ksa .