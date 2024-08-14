# Script to build your protobuf, go binaries, and docker images here
# Generate Go binaries
rm bin/*
GOOS=linux GOARCH=amd64 go build -o ./bin/master cmd/master/*.go
GOOS=linux GOARCH=amd64 go build -o ./bin/worker cmd/worker/main.go
cp ./cmd/worker/mapper/map.py ./bin/.
cp ./cmd/worker/reducer/reduce.py ./bin/.

# sudo docker run --name=master greeter-master tail -f /dev/null
# sudo docker run --name=worker greeter-worker tail -f /dev/null


# sudo docker kill $(sudo docker ps -aq)
# sudo docker rm $(sudo docker ps -aq)
