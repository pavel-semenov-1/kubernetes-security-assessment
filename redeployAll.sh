#!/bin/bash

helm -n demo-users upgrade --install demo-users k8s/demo-users
helm -n demo-subjects upgrade --install demo-subjects k8s/demo-subjects