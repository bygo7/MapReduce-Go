#!/bin/bash

# Install Kubectl
sudo curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
echo "$(<kubectl.sha256) kubectl" | sha256sum --check

# Install Protobuf
brew install protobuf
protoc --version

# Install Kind
export PATH=$PATH:/usr/local/go/bin
brew install kind
export PATH="$PATH:$(go env GOPATH)/bin"

# Install Helm
brew install helm

# Install Go GRPC
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

# Setup go module
go mod init mapreduce
go mod tidy

# Clean up
sudo rm kubectl
sudo rm kubectl.sha256
