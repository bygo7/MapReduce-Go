import sys

# Arg 1: Docker repository
# Arg 2: Docker image tag
DOCKER_REPOSITORY = sys.argv[1]
DOCKER_TAG = sys.argv[2]
WORKER_SIZE = int(sys.argv[3])
MASTER_FAIL_NODE = int(sys.argv[4])
MASTER_FAIL_POINT = int(sys.argv[5])
WORKER_FAIL_NODE = int(sys.argv[6])
WORKER_FAIL_POINT = int(sys.argv[7])
EPOCH_TIME = int(sys.argv[8])

if DOCKER_REPOSITORY == "":
    DOCKER_REPOSITORY = "mapreduce-"
if DOCKER_TAG == "":
    DOCKER_TAG = "g"

master_service = """
apiVersion: v1
kind: Service
metadata:
    name: master-service
spec:
    selector:
        app: mapreduce-master
    clusterIP: None
    ports:
        - protocol: TCP
          port: 8080
          targetPort: 8080
---
"""
master_service += """
apiVersion: v1
kind: Service
metadata:
    name: master-ingress
    annotations:
        service.beta.kubernetes.io/azure-load-balancer-resource-group: MC_MapReduce_group_MapReduce_eastus
        service.beta.kubernetes.io/azure-pip-name: MapReduceGatewayIP
        service.beta.kubernetes.io/azure-dns-label-name: mapreducet20
spec:
    type: LoadBalancer
    ports:
        - port: 80
          targetPort: 8090
    selector:
        app: mapreduce-master
---
"""
worker_services = ""

for i in range(0, WORKER_SIZE):
    worker_service = f"""
apiVersion: v1
kind: Service
metadata:
    name: worker-service-{i}
spec:
    selector:
        app: mapreduce-worker-{i}
    ports:
        - protocol: TCP
          port: 8080
          targetPort: 8080
---
"""
    worker_services += worker_service

master_pod = f"""
apiVersion: apps/v1
kind: Deployment
metadata:
    name: master-pod
spec:
    replicas: 2
    selector:
        matchLabels:
            app: mapreduce-master
    template:
        metadata:
            labels:
                app: mapreduce-master
        spec:
            containers:
                - name: master-container
                  image: {DOCKER_REPOSITORY}master:{DOCKER_TAG}
                  imagePullPolicy: IfNotPresent
                  args: ["-name=master", "-logtostderr=true", "-failureNode={MASTER_FAIL_NODE}", "-failurePoint={MASTER_FAIL_POINT}", "-launchTime={EPOCH_TIME}"]
                  ports:
                  - containerPort: 8080
                  - containerPort: 8090
                  env:
                  - name: POD_IP
                    valueFrom:
                        fieldRef:
                            fieldPath: status.podIP
                  livenessProbe:
                        httpGet:
                            path: /healthz
                            port: 8090
                        initialDelaySeconds: 5
                        periodSeconds: 5

---
"""
worker_pods = """"""
for i in range(0, WORKER_SIZE):
    failurePoint = "0"
    if i == WORKER_FAIL_NODE:
        failurePoint = WORKER_FAIL_POINT
    worker_pod = f"""
apiVersion: apps/v1
kind: Deployment
metadata:
    name: worker-pod-{i}
spec:
    replicas: 1
    selector:
        matchLabels:
            app: mapreduce-worker-{i}
    template:
        metadata:
            labels:
                app: mapreduce-worker-{i}
        spec:
            containers:
            - name: worker-container
              image: {DOCKER_REPOSITORY}worker:{DOCKER_TAG}
              imagePullPolicy: IfNotPresent
              args: ["-name=worker-{i}", "-serviceName=worker-service-{i}", "-logtostderr=true", "-failurePoint={failurePoint}"]
              ports:
              - containerPort: 8080
---
"""
    worker_pods += worker_pod

s = master_service + worker_services + master_pod + worker_pods
f = open("kubernetes.yaml", "w")
f.write(s)
f.close()
