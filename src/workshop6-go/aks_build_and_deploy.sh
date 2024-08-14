#!/bin/bash

# Workflow:
# 1. Build Go binaries
# 2. Build Docker images, set tag as the current time
# 3. Push Docker images to ACR
# 4. Generate new Kubernetes yaml file
# 5. Deploy to AKS

# Constants
ACR_NAME=mapreducet20registry.azurecr.io/mapreduce
TAG=$(date '+%Y%m%d-%H%M%S')
MASTER_DOCKER_IMAGE=$ACR_NAME"/master:"$TAG
WORKER_DOCKER_IMAGE=$ACR_NAME"/worker:"$TAG
EPOCH_TIME=$(date +%s)

# App configuration constants
WORKER_SIZE=7
MASTER_FAIL_NODE=-1
MASTER_FAIL_POINT=21
WORKER_FAIL_NODE=-1
WORKER_FAIL_POINT=3

# 0. Login to ACR
# az login
az acr login --name mapreducet20registry

# 1. Build Go binaries
./build.sh

# 2. Build Docker images, set tag as the current time
docker build \
    --build-arg AZURE_STORAGE_ACCOUNT_NAME=$AZURE_STORAGE_ACCOUNT_NAME \
    --build-arg AZURE_STORAGE_ACCOUNT_KEY=$AZURE_STORAGE_ACCOUNT_KEY \
    -t $MASTER_DOCKER_IMAGE -f docker/Dockerfile.master .
docker build \
    --build-arg AZURE_STORAGE_ACCOUNT_NAME=$AZURE_STORAGE_ACCOUNT_NAME \
    --build-arg AZURE_STORAGE_ACCOUNT_KEY=$AZURE_STORAGE_ACCOUNT_KEY \
    -t $WORKER_DOCKER_IMAGE -f docker/Dockerfile.worker .

# 3. Push Docker images to ACR
docker push $MASTER_DOCKER_IMAGE
docker push $WORKER_DOCKER_IMAGE

# 4. Generate new Kubernetes yaml file
python3 gen_kubernetes_yaml.py $ACR_NAME/ $TAG $WORKER_SIZE $MASTER_FAIL_NODE $MASTER_FAIL_POINT $WORKER_FAIL_NODE $WORKER_FAIL_POINT $EPOCH_TIME

# 5. Deploy to AKS
kubectl config set-context MapReduce
kubectl create ns mapreduce
helm repo add bitnami https://charts.bitnami.com/bitnami
# helm install etcd bitnami/etcd --set auth.rbac.create=false --namespace mapreduce
helm install etcd bitnami/etcd --set auth.rbac.create=false,readinessProbe.enabled=false,livenessProbe.enabled=false,startupProbe.enabled=false --namespace mapreduce
kubectl delete pod etcd-client --namespace mapreduce
kubectl run etcd-client --restart='Never' --image docker.io/bitnami/etcd:3.5.10-debian-11-r1 --env ETCDCTL_ENDPOINTS="etcd.mapreduce.svc.cluster.local:2379" --namespace mapreduce --command -- sleep infinity
kubectl wait --for=condition=ready pod/etcd-client --timeout=120s --namespace mapreduce
kubectl apply -f kubernetes.yaml -n mapreduce
