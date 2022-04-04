#!/usr/bin/env bash
eval $(minikube docker-env)
skaffold dev #--no-prune=false --cache-artifacts=false
