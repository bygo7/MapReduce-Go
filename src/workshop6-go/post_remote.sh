#!/bin/bash
# Get first argument
IP=$1
# curl -X POST -H "Content-Type: application/json" -d '{"containerName":"gutenberg", "m": 20, "r": 5}' http://$IP:80/word-count
curl -X POST -H "Content-Type: application/json" -d '{"containerName":"test-1gb", "m": 100, "r": 20}' http://$IP:80/word-count

# mapreducet20.eastus.cloudapp.azure.com
# 20.119.125.60
