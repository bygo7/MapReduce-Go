#!/bin/bash
# sudo -E required for local environment variables

# Setup kind
kind create cluster

MASTER_DOCKER_IMAGE=mapreduce-master:g
WORKER_DOCKER_IMAGE=mapreduce-worker:g

# Build Docker images
docker build \
    --build-arg AZURE_STORAGE_ACCOUNT_NAME=$AZURE_STORAGE_ACCOUNT_NAME \
    --build-arg AZURE_STORAGE_ACCOUNT_KEY=$AZURE_STORAGE_ACCOUNT_KEY \
    -t $MASTER_DOCKER_IMAGE -f docker/Dockerfile.master .
docker build \
    --build-arg AZURE_STORAGE_ACCOUNT_NAME=$AZURE_STORAGE_ACCOUNT_NAME \
    --build-arg AZURE_STORAGE_ACCOUNT_KEY=$AZURE_STORAGE_ACCOUNT_KEY \
    -t $WORKER_DOCKER_IMAGE -f docker/Dockerfile.worker .

# load docker images
kind load docker-image $MASTER_DOCKER_IMAGE
kind load docker-image $WORKER_DOCKER_IMAGE

# Build kubernetes
kubectl create ns mapreduce

# helm repo add bitnami https://charts.bitnami.com/bitnami
# helm install etcd bitnami/etcd --set auth.rbac.create=false,readinessProbe.enabled=false,livenessProbe.enabled=false,startupProbe.enabled=false --namespace helloworld

# # generate kubernetes yaml file
# python3 gen_kubernetes_yaml.py

# Deploy
kubectl -n mapreduce apply -f kubernetes.yaml

kubectl config set-context --current --namespace=mapreduce

# kubectl get all -n mapreduce 

# kubectl logs master-pod-1-768dd69b66-9mbqs
# sudo kubectl describe pod etcd-0 -n helloworld
