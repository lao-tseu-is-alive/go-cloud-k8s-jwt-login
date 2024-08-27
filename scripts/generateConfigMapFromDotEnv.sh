#!/bin/bash
kubectl create configmap app-config-go-cloud-k8s-jwt-login -n go-testing --from-env-file=.env --output yaml --dry-run=client >deployments/go-testing/configmap-local.yaml
